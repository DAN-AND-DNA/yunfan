package pkg

import (
	"io"
	"net"
	"time"

	"testing"
	pkg_rpc "yunfan/pkg/rpc"

	"snk.git.node1/yunfan/arpc"
	"snk.git.node1/yunfan/arpc/jsonrpc"
)

type Dan struct {
}

func (this *Dan) Ping(rwc io.ReadWriteCloser, args *string, reply *string) error {
	*reply = *args
	return nil
}

func Test_pkg_json_rpc_call_easy(t *testing.T) {
	cli, srv := net.Pipe()
	server := pkg_rpc.New_server()

	err := server.Register(&Dan{})
	if err != nil {
		panic(err)
	}

	go server.Test_json(srv)

	ptr_client := jsonrpc.NewClient(cli)
	defer ptr_client.Close()

	var reply string
	if err := ptr_client.Call("Dan.Ping", "ping", &reply); err != nil {
		panic(err)
	}

	if reply != "ping" {
		panic(reply)
	}
}

type Yang struct {
}

type C_args struct {
	Bytes_vals []byte
	Carrier    map[string]string
}

func (this *Yang) Ping(rwc io.ReadWriteCloser, args *C_args, reply *string) error {
	*reply = args.Carrier["ping"]
	return nil
}
func (this *Yang) Ping_raw(rwc io.ReadWriteCloser, args map[string]string, reply *string) error {
	*reply = args["ping"]
	return nil
}
func (this *Yang) Ping_bytes(rwc io.ReadWriteCloser, args *C_args, reply *string) error {
	*reply = string(args.Bytes_vals)
	return nil
}

func Test_pkg_json_rpc_call_map(t *testing.T) {
	cli, srv := net.Pipe()
	server := pkg_rpc.New_server()

	err := server.Register(&Yang{})
	if err != nil {
		panic(err)
	}

	go server.Test_json(srv)

	ptr_client := jsonrpc.NewClient(cli)
	defer ptr_client.Close()

	var reply string
	m := make(map[string]string)
	m["ping"] = "pong"
	if err := ptr_client.Call("Yang.Ping", &C_args{Bytes_vals: []byte(`dadadaddcccccccc`), Carrier: m}, &reply); err != nil {
		panic(err)
	}

	if reply != "pong" {
		panic(reply)
	}

	var reply1 string
	if err := ptr_client.Call("Yang.Ping_raw", m, &reply1); err != nil {
		panic(err)
	}

	if reply1 != "pong" {
		panic(reply1)
	}

}

func Test_pkg_json_rpc_call_bytes(t *testing.T) {
	cli, srv := net.Pipe()
	server := pkg_rpc.New_server()

	err := server.Register(&Yang{})
	if err != nil {
		panic(err)
	}

	go server.Test_json(srv)

	ptr_client := jsonrpc.NewClient(cli)
	defer ptr_client.Close()

	var reply string
	if err := ptr_client.Call("Yang.Ping_bytes", &C_args{Bytes_vals: []byte(`danyangchen`)}, &reply); err != nil {
		panic(err)
	}

	if reply != "danyangchen" {
		panic(reply)
	}

}

func Test_pkg_rpc_pool(t *testing.T) {
	c_conn, s_conn := net.Pipe()
	server := pkg_rpc.New_server()

	err := server.Register(&Dan{})
	if err != nil {
		panic(err)
	}

	go server.Test_json(s_conn)
	pool, err := pkg_rpc.New_json_rpc_client_pool(1, 1,
		func() (net.Conn, error) {
			return c_conn, nil
		},

		func(rpc_client *arpc.Client) bool {
			var reply string
			if err := rpc_client.Call("Dan.Ping", "ping", &reply); err != nil {
				panic(err)
			}

			if reply != "ping" {
				panic(reply)
			}

			return true
		},
		1,
	)

	if err != nil {
		panic(err)
	}

	/*
		for i := 0; i < 5; i++ {
			client, err := pool.Get()
			if err != nil {
				panic(err)
			}
			client.Close()
		}
	*/

	_ = pool

	time.Sleep(3 * time.Second)
}
