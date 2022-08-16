package conn_pool

import (
	"errors"
	"log"
	"net"
	"sync"
)

var (
	Err_the_pool_is_closed                      = errors.New("error: the pool is closed")
	Err_bad_pre_allocated_size                  = errors.New("error: bad pre allocated size")
	Err_bad_max_size                            = errors.New("error: bad max size")
	Err_pre_allocated_size_bigger_than_max_size = errors.New("error: pre allocated size bigger than max size")
	Err_pool_no_init                            = errors.New("errors: the pool is no init")
)

type Conn_pool interface {
	Get() (net.Conn, error)

	Close()

	Len() int

	Keep_alive() bool

	put(net.Conn)
}

type Conn_factory func() (net.Conn, error)
type Keep_alive_callback func() bool

func New(pre_allocated_size, max_size int, conn_factory Conn_factory, keep_alive_callback Keep_alive_callback, raw_is_mock ...bool) (Conn_pool, error) {
	if pre_allocated_size < 0 {
		return nil, Err_bad_pre_allocated_size
	}

	if max_size <= 0 {
		return nil, Err_bad_max_size
	}

	if pre_allocated_size > max_size {
		return nil, Err_pre_allocated_size_bigger_than_max_size
	}

	is_mock := false
	if len(raw_is_mock) == 0 {
		is_mock = raw_is_mock[0]
	}

	if is_mock {
		mock_pool := &Mock_pool{}
		return mock_pool, nil
	}

	channel_pool := &Channel_pool{
		is_init:             true,
		conns:               make(chan net.Conn, max_size),
		conn_factory:        conn_factory,
		keep_alive_callback: keep_alive_callback,
	}

	for i := 0; i < pre_allocated_size; i++ {
		if conn, err := conn_factory(); err != nil {
			log.Printf("error: pre allocate connection failed where init a pool: %s\n", err.Error())
			break
		} else {
			channel_pool.conns <- conn
			channel_pool.current_conn_size++
		}
	}

	return channel_pool, nil
}

type Channel_pool struct {
	is_init             bool
	conn_factory        Conn_factory
	conns               chan net.Conn
	current_conn_size   int
	m_conns             sync.RWMutex
	keep_alive_callback Keep_alive_callback
}

func (this *Channel_pool) Get() (net.Conn, error) {
	select {
	case raw_conn := <-this.conns:
		if raw_conn == nil {
			return nil, Err_the_pool_is_closed
		}

		return this.new_conn(raw_conn), nil
	default:
		raw_conn, err := this.conn_factory()
		if err != nil {
			return nil, err
		}

		return this.new_conn(raw_conn), nil
	}

	return nil, nil
}

func (this *Channel_pool) Close() {
	this.m_conns.Lock()
	if this.is_init == false {
		this.m_conns.Unlock()
		return
	}

	this.is_init = false
	this.m_conns.Unlock()

	close(this.conns)
	for conn := range this.conns {
		conn.Close()
	}
}

func (this *Channel_pool) Len() int {
	this.m_conns.RLock()
	defer this.m_conns.RUnlock()

	if this.is_init == false {
		return 0
	}

	return this.current_conn_size
}

func (this *Channel_pool) Keep_alive() bool {
	this.m_conns.RLock()
	defer this.m_conns.RUnlock()
	if this.is_init == false {
		return false
	}

	return this.keep_alive_callback()
}

func (this *Channel_pool) new_conn(raw_conn net.Conn) net.Conn {
	pconn := &Pool_conn{is_init: true, Conn: raw_conn, owner_pool: this, is_bad: false}
	return pconn
}

func (this *Channel_pool) put(raw_conn net.Conn) {
	if raw_conn == nil {
		return
	}

	this.m_conns.RLock()
	defer this.m_conns.RUnlock()

	if this.is_init == false {
		raw_conn.Close()
		return
	}

	select {
	case this.conns <- raw_conn:
		return
	default:
		raw_conn.Close()
	}

}

type Mock_pool struct {
}

func (this *Mock_pool) Get() (net.Conn, error) {
	return nil, nil
}

func (this *Mock_pool) Close() {
}

func (this *Mock_pool) Len() int {
	return 0
}

func (this *Mock_pool) Keep_alive() bool {
	return true
}

func (this *Mock_pool) put(net.Conn) {
}
