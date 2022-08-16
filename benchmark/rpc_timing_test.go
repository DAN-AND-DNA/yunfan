package timing_test

import (
	"bytes"
	"encoding/gob"

	"io"
	"testing"
	pkg_rpc "yunfan/pkg/rpc"
	pkg_test_tool "yunfan/pkg/test_tool"

	json "github.com/json-iterator/go"
	"snk.git.node1/yunfan/arpc"
)

type Dan struct {
}

func (this *Dan) Ping(rwc io.ReadWriteCloser, args *string, reply *string) error {
	*reply = *args
	return nil
}

type Args struct {
	V string
}

type Reply struct {
	V string
}

type DanPingResp struct {
	Id     interface{} `json:"id"`
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
}

func Test_gob_rpc_call(_ *testing.T) {

	r := &arpc.Request{}
	r.Seq = 37
	r.ServiceMethod = "Dan.Ping"

	var enc_buf bytes.Buffer
	enc := gob.NewEncoder(&enc_buf)

	var args string = "ping"
	var reply string
	enc.Encode(r)
	enc.Encode(&args)

	dump_req := enc_buf.Bytes()

	conn := new(pkg_test_tool.Test_conn)
	conn.Input_req(dump_req)

	server := pkg_rpc.New_server()
	d := &Dan{}
	if err := server.Register(d); err != nil {
		panic(err)
	}
	server.Test(conn)

	buffer := make([]byte, 200)
	n, err := conn.Output_resp(buffer)
	if err != nil {
		panic(err)
	}

	buffer = buffer[:n]
	dec_buf := bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(dec_buf)

	r1 := &arpc.Request{}
	dec.Decode(r1)
	dec.Decode(&reply)

	if reply != "ping" {
		panic("bad reply")
	}

	conn.Reset()
}

func Test_json_rpc_call(_ *testing.T) {
	dump_req := []byte(`{"method": "Dan.Ping", "id": "7", "params":["ping"]}`)
	conn := new(pkg_test_tool.Test_conn)
	conn.Input_req(dump_req)

	server := pkg_rpc.New_server()
	d := &Dan{}
	if err := server.Register(d); err != nil {
		panic(err)
	}
	server.Test_json(conn)

	buffer := make([]byte, 400)
	if n, err := conn.Output_resp(buffer); err != nil {
		panic(err)
	} else {
		buffer = buffer[:n]
	}

	dec_buf := bytes.NewBuffer(buffer)
	dec := json.NewDecoder(dec_buf)
	reply := DanPingResp{}
	if err := dec.Decode(&reply); err != nil {
		panic(err)
	}

	if reply.Result != "ping" {
		panic("bad pong")
	}

	conn.Reset()
}

func Benchmark_server_handle_json_rpc_call(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		dump_req := []byte(`{"method": "Dan.Ping", "id": "7", "params":["ping"]}`)

		conn := new(pkg_test_tool.Test_conn)

		server := pkg_rpc.New_server()
		d := &Dan{}
		if err := server.Register(d); err != nil {
			panic(err)
		}

		for pb.Next() {
			conn.Input_req(dump_req)
			server.Test_json(conn)
			conn.Reset()
		}
	})
}
func Benchmark_server_handle_gob_rpc_call(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		r := &arpc.Request{}
		r.Seq = 37
		r.ServiceMethod = "Dan.Ping"

		var enc_buf bytes.Buffer
		enc := gob.NewEncoder(&enc_buf)

		var args string = "ping"
		enc.Encode(r)
		enc.Encode(&args)

		dump_req := enc_buf.Bytes()

		conn := new(pkg_test_tool.Test_conn)
		server := pkg_rpc.New_server()
		d := &Dan{}
		if err := server.Register(d); err != nil {
			panic(err)
		}

		for pb.Next() {
			conn.Input_req(dump_req)
			server.Test(conn)
			conn.Reset()
		}
	})
}
