package rpc

import (
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"snk.git.node1/yunfan/arpc"
	"snk.git.node1/yunfan/arpc/jsonrpc"
)

var (
	Err_the_pool_is_closed                      = errors.New("error: the pool is closed")
	Err_bad_pre_allocated_size                  = errors.New("error: bad pre allocated size")
	Err_bad_max_size                            = errors.New("error: bad max size")
	Err_pre_allocated_size_bigger_than_max_size = errors.New("error: pre allocated size bigger than max size")
	Err_pool_no_init                            = errors.New("error: the pool is no init")
)

type Std_rpc_client_pool interface {
	Get() (*Json_rpc_client, error)

	Close()

	Len() int

	Keep_alive(time.Duration)

	put(arpc_client *arpc.Client)
}

type Conn_factory func() (net.Conn, error)
type Keep_alive_callback func(rpc_client *arpc.Client) bool

func New_json_rpc_client_pool(pre_allocated_size, max_size int, conn_factory Conn_factory, keep_alive_callback Keep_alive_callback, keep_alive_interval int) (Std_rpc_client_pool, error) {
	if pre_allocated_size < 0 {
		return nil, Err_bad_pre_allocated_size
	}

	if max_size <= 0 {
		return nil, Err_bad_max_size
	}

	if pre_allocated_size > max_size {
		return nil, Err_pre_allocated_size_bigger_than_max_size
	}

	json_rpc_client_pool := &Json_rpc_client_pool{
		is_init:             true,
		conn_factory:        conn_factory,
		clients:             make(chan *arpc.Client, max_size),
		keep_alive_callback: keep_alive_callback,
	}

	for i := 0; i < pre_allocated_size; i++ {
		if conn, err := conn_factory(); err != nil {
			log.Printf("error: pre allocate connection failed where init a pool: %s\n", err.Error())
			break
		} else if conn != nil {
			client := jsonrpc.NewClient(conn)
			id := *(*uint64)(unsafe.Pointer(client))
			json_rpc_client_pool.clients <- client
			json_rpc_client_pool.current_client_size++
			json_rpc_client_pool.client_stores.Store(id, &Json_rpc_client_store{})
		}
	}

	go json_rpc_client_pool.Keep_alive((time.Duration)(keep_alive_interval) * time.Second)

	return json_rpc_client_pool, nil

}

type Json_rpc_client_pool struct {
	is_init             bool
	conn_factory        Conn_factory
	clients             chan *arpc.Client
	current_client_size int
	keep_alive_callback Keep_alive_callback
	m_clients           sync.RWMutex
	client_stores       sync.Map
}

func (this *Json_rpc_client_pool) new_json_rpc_client(rpc_client *arpc.Client) *Json_rpc_client {
	return &Json_rpc_client{
		Client:     rpc_client,
		is_init:    true,
		owner_pool: this,
		is_closed:  0,
	}
}

func (this *Json_rpc_client_pool) Get() (*Json_rpc_client, error) {
	select {
	case client := <-this.clients:
		if client == nil {
			return nil, Err_the_pool_is_closed
		}

		id := *(*uint64)(unsafe.Pointer(client))
		raw_val, ok := this.client_stores.LoadAndDelete(id)
		var val *Json_rpc_client_store
		if ok {
			val = raw_val.(*Json_rpc_client_store)
			val.get_times++
			if val.get_times >= 3 {
				val.get_times = 0
				this.keep_alive_callback(client)
			}

			this.client_stores.Store(id, val)
		}

		return this.new_json_rpc_client(client), nil
	default:
		conn, err := this.conn_factory()
		if err != nil {
			return nil, err
		}

		if conn == nil {
			// this is mock
			return this.new_json_rpc_client(nil), nil
		}

		client := jsonrpc.NewClient(conn)
		id := *(*uint64)(unsafe.Pointer(client))
		this.client_stores.Store(id, &Json_rpc_client_store{})

		return this.new_json_rpc_client(client), nil
	}
}

func (this *Json_rpc_client_pool) Close() {
	this.m_clients.Lock()
	if this.is_init == false {
		this.m_clients.Unlock()
		return
	}

	this.is_init = false
	this.m_clients.Unlock()

	close(this.clients)
	for client := range this.clients {

		id := *(*uint64)(unsafe.Pointer(client))
		this.client_stores.Delete(id)
		client.Close()
	}
}

func (this *Json_rpc_client_pool) Len() int {
	this.m_clients.RLock()
	defer this.m_clients.RUnlock()

	if this.is_init == false {
		return 0
	}

	return this.current_client_size

}

func (this *Json_rpc_client_pool) Keep_alive(interval time.Duration) {
	for {
		time.Sleep(interval)

		this.m_clients.RLock()
		if this.is_init == false {
			this.m_clients.RUnlock()
			continue
		}
		this.m_clients.RUnlock()

		c_len := this.Len()

		for i := 0; i < c_len; i++ {
			select {
			case client := <-this.clients:
				if !this.keep_alive_callback(client) {
					id := *(*uint64)(unsafe.Pointer(client))
					this.client_stores.Delete(id)
					client.Close()
					log.Println("disconnect from server")
				} else {
					select {
					case this.clients <- client:
						continue
					default:
						id := *(*uint64)(unsafe.Pointer(client))
						this.client_stores.Delete(id)
						client.Close()
					}
				}
			default:
				continue
			}
		}
	}
}

func (this *Json_rpc_client_pool) put(raw_client *arpc.Client) {
	if raw_client == nil {
		return
	}

	this.m_clients.RLock()
	defer this.m_clients.RUnlock()

	if this.is_init == false {
		raw_client.Close()
		return
	}

	select {
	case this.clients <- raw_client:
		return
	default:
		id := *(*uint64)(unsafe.Pointer(raw_client))
		this.client_stores.Delete(id)
		raw_client.Close()
	}
}

type Json_rpc_client struct {
	*arpc.Client
	is_init    bool
	is_closed  uint32
	is_bad     bool
	owner_pool *Json_rpc_client_pool
}

type Json_rpc_client_store struct {
	get_times int
}

func (this *Json_rpc_client) Close() {
	if !this.is_init || this.Client == nil {
		return
	}

	if swapped := atomic.CompareAndSwapUint32(&this.is_closed, 0, 1); !swapped {
		return
	}

	if this.is_bad {
		id := *(*uint64)(unsafe.Pointer(this.Client))
		this.owner_pool.client_stores.Delete(id)

		this.Client.Close()
		return
	}

	this.owner_pool.put(this.Client)
}

func (this *Json_rpc_client) Set_bad() {
	if !this.is_init {
		return
	}

	this.is_bad = true
}
