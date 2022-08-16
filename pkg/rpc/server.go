package rpc

import (
	"crypto/sha1"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/crypto/pbkdf2"

	"io"

	krpc "snk.git.node1/yunfan/arpc"
	"snk.git.node1/yunfan/arpc/jsonrpc"
	//krpc "net/rpc"
)

var (
	default_server = New_server()

	ERR_IS_RUNNING        = errors.New("rpc server is running")
	ERR_BAD_PROTOCOL_TYPE = errors.New("bad protocol type")
	ERR_LISTEN_FAILED     = errors.New("listen failed: already try 3 times!!!")
	ERR_NEED_KEY_SALT     = errors.New("need key and salt for udp")
)

type Rpc_server struct {
	is_running   bool
	krpc_server  *krpc.Server
	listener     net.Listener
	m_server     sync.Mutex
	m_err_server sync.Mutex
	err_server   error
}

func New_server() *Rpc_server {
	return &Rpc_server{
		is_running:  false,
		krpc_server: krpc.NewServer(),
	}
}

type Local_service int

func (this *Local_service) Are_u_ok(rwc io.ReadWriteCloser, arg *int, reply *int) error {
	*reply = *arg
	return nil
}

type Args struct {
	Val int
}

type Reply struct {
	Val int
}

type Local_service_json int

func (this *Local_service_json) Are_u_ok_json(rwc io.ReadWriteCloser, args *Args, reply *Reply) error {
	reply.Val = args.Val
	return nil
}

func (this *Rpc_server) get_error() error {
	this.m_err_server.Lock()
	defer this.m_err_server.Unlock()

	if this.err_server == nil {
		return nil
	}

	return errors.New(this.err_server.Error())
}

func (this *Rpc_server) listen(protocol_type, app_protocol_type, port, key, salt string, idle_s int) error {
	if this.is_running {
		return ERR_IS_RUNNING
	}

	if protocol_type == "" || port == "" {
		return nil
	}

	switch protocol_type {
	case "tcp":
		ptr_local_service := new(Local_service)
		ptr_local_service_json := new(Local_service_json)

		err := this.krpc_server.Register(ptr_local_service)
		if err != nil {
			return err
		}
		err = this.krpc_server.Register(ptr_local_service_json)
		if err != nil {
			return err
		}

		this.listener, err = net.Listen("tcp", "0.0.0.0:"+port)
		if err != nil {
			return err
		}

		go func() {
			for {
				conn, err := this.listener.Accept()
				if err != nil {
					this.m_err_server.Lock()
					this.err_server = err
					this.m_err_server.Unlock()
					return
				}

				tcp_conn := conn.(*net.TCPConn)
				tcp_conn.SetNoDelay(true)
				tcp_conn.SetReadBuffer(4 * 1024 * 1024)
				tcp_conn.SetReadBuffer(4 * 1024 * 1024)
				if idle_s <= 0 {
					idle_s = 5
				}

				log.Printf("new tcp conn: %ds life", idle_s)
				tcp_conn.SetDeadline(time.Now().Add(time.Duration(idle_s) * time.Second))
				if app_protocol_type == "json" {
					go this.krpc_server.ServeCodec(jsonrpc.NewServerCodec(conn))
				} else {
					go this.krpc_server.ServeConn(conn)
				}
			}
		}()

		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			tcp_ccon, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 1*time.Second)
			if err != nil {
				log.Println(err)
				continue
			}

			tcp_ccon.SetDeadline(time.Now().Add(5 * time.Second))

			if app_protocol_type == "json" {
				ptr_client := jsonrpc.NewClient(tcp_ccon)
				var reply Reply
				if err := ptr_client.Call("Local_service_json.Are_u_ok_json", &Args{Val: 7}, &reply); err == nil && reply.Val == 7 {
					this.is_running = true
					ptr_client.Close()
					break
				} else {
					log.Println(err)
				}

				ptr_client.Close()

			} else {
				ptr_client := krpc.NewClient(tcp_ccon)
				var reply int
				if err := ptr_client.Call("Local_service.Are_u_ok", 7, &reply); err == nil && reply == 7 {
					this.is_running = true
					ptr_client.Close()
					break
				}

				ptr_client.Close()
			}

		}

	case "udp":
		if len(key) == 0 || len(salt) == 0 {
			return ERR_NEED_KEY_SALT
		}

		pass := pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha1.New)
		block_pass, err := kcp.NewAESBlockCrypt(pass[:16]) // AES128
		if err != nil {
			return err
		}

		ptr_local_service := new(Local_service)
		err = this.krpc_server.Register(ptr_local_service)
		if err != nil {
			return err
		}

		this.listener, err = kcp.ListenWithOptions("0.0.0.0:"+port, block_pass, 10, 3)
		if err != nil {
			return err
		}

		kcp_listener := this.listener.(*kcp.Listener)

		kcp_listener.SetReadBuffer(4 * 1024 * 1024)
		kcp_listener.SetWriteBuffer(4 * 1024 * 1024)
		kcp_listener.SetDSCP(46)

		go func() {
			for {
				conn, err := this.listener.Accept()
				if err != nil {
					this.m_err_server.Lock()
					this.err_server = err
					this.m_err_server.Unlock()
					return
				}

				kcp_conn := conn.(*kcp.UDPSession)
				kcp_conn.SetStreamMode(true)
				kcp_conn.SetWindowSize(4096, 4096)
				kcp_conn.SetNoDelay(1, 10, 2, 1)
				kcp_conn.SetDSCP(46)
				kcp_conn.SetMtu(1400)
				kcp_conn.SetACKNoDelay(false)
				kcp_conn.SetWriteDelay(false)
				if idle_s <= 0 {
					idle_s = 3
				}

				kcp_conn.SetDeadline(time.Now().Add(time.Duration(idle_s) * time.Second))
				log.Printf("new kcp conn : %ds life", idle_s)
				if app_protocol_type == "json" {
					go this.krpc_server.ServeCodec(jsonrpc.NewServerCodec(conn))
				} else {
					go this.krpc_server.ServeConn(conn)
				}
			}
		}()

		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			kcp_ccon, err := kcp.DialWithOptions("127.0.0.1:"+port, block_pass, 10, 3)
			if err != nil {
				continue
			}

			kcp_ccon.SetDeadline(time.Now().Add(2 * time.Second))
			ptr_client := krpc.NewClient(kcp_ccon)
			var reply int
			if err := ptr_client.Call("Local_service.Are_u_ok", 7, &reply); err == nil && reply == 7 {
				this.is_running = true
				ptr_client.Close()
				break
			}

			ptr_client.Close()
		}

	default:
		return ERR_BAD_PROTOCOL_TYPE
	}

	if !this.is_running {
		if err := this.get_error(); err != nil {
			return err
		} else {
			return ERR_LISTEN_FAILED
		}
	}

	log.Println("rpc server: ok")
	return nil
}

func (this *Rpc_server) Register(receiver interface{}) error {
	if this.is_running {
		return ERR_IS_RUNNING
	}

	return this.krpc_server.Register(receiver)
}

func (this *Rpc_server) Register_by_name(receiver interface{}, service_name string) error {
	if this.is_running {
		return ERR_IS_RUNNING
	}

	return this.krpc_server.RegisterName(service_name, receiver)
}

func (this *Rpc_server) shutdown(deadline int) {
	if !this.is_running {
		return
	}

	go func() {
		this.listener.Close()
	}()

	log.Println("rpc server: shutdown")
	time.Sleep((time.Duration)(deadline) * time.Second)
}

func (this *Rpc_server) Test(conn io.ReadWriteCloser) {
	this.krpc_server.ServeConn(conn)
}

func (this *Rpc_server) Test_json(conn io.ReadWriteCloser) {
	this.krpc_server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func Listen(protocol_type, app_protocol_type, port, key, salt string, idle int) error {
	return default_server.listen(protocol_type, app_protocol_type, port, key, salt, idle)
}

func Register_by_name(receiver interface{}, service_name string) error {
	return default_server.Register_by_name(receiver, service_name)
}

func Register(receiver interface{}) error {
	return default_server.Register(receiver)
}

func Shutdown(deadline int) {
	default_server.shutdown(deadline)
}
