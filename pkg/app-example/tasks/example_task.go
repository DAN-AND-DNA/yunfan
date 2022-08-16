package tasks

import (
	"context"
	"net"
	"time"
	dev_log "yunfan/pkg/log"
	"yunfan/pkg/rpc"

	"snk.git.node1/yunfan/arpc"
	krpc "snk.git.node1/yunfan/arpc"
)

type Example_task struct {
	num int
}

func New_example_task() *Example_task {
	return &Example_task{num: 0}
}

type A struct {
	Age int `json:"age"`
}

func (this *Example_task) Do_local_log(current_execute_num uint64) error {
	str_json_val := "json_val"
	dev_log.Debug("i", this.num, "msg", "this is a debug", str_json_val, A{Age: 73})
	this.num++
	dev_log.Info("i", this.num, "msg", "this is a info", str_json_val, A{Age: 73})
	this.num++
	dev_log.Warn("i", this.num, "msg", "this is a warn", str_json_val, A{Age: 73})
	this.num++
	dev_log.Error("i", this.num, "msg", "this is a error", str_json_val, A{Age: 73})
	this.num++
	return nil
}

func (this *Example_task) Do_rpc_tracer(current_exaute_num uint64) error {
	rpc_span := rpc.New_rpc_span()
	ctx, err := rpc_span.Start(context.Background(), "do-rpc-tracer")
	if err != nil {
		return err
	}
	defer rpc_span.Finish()

	tcp_ccon, err := net.DialTimeout("tcp", "127.0.0.1:3777", 2*time.Second)
	if err != nil {
		return err
	}

	tcp_ccon.SetDeadline(time.Now().Add(3 * time.Second))
	ptr_client := krpc.NewClient(tcp_ccon)
	defer ptr_client.Close()

	if err := f1(ptr_client, ctx); err != nil {
		return err
	}

	if err := f2(ptr_client, ctx); err != nil {
		return err
	}

	if err := f3(ctx); err != nil {
		return err
	}

	return nil
}

type Ping_args struct {
	Carrier map[string]string
}

func f1(ptr_client *arpc.Client, ctx context.Context) error {
	rpc_span := rpc.New_rpc_span()
	ctx, err := rpc_span.Start(ctx, "f1")
	if err != nil {
		return err
	}
	defer rpc_span.Finish()

	var reply string
	args := &Ping_args{Carrier: rpc_span.Carrier}

	if err := ptr_client.Call("v1.api.a_service.Ping", args, &reply); err != nil {
		return err
	}

	dev_log.Info("msg", reply)

	if err := f11(ctx); err != nil {
		return err
	}

	return nil
}

func f11(ctx context.Context) error {
	rpc_span := rpc.New_rpc_span()
	if _, err := rpc_span.Start(ctx, "f11"); err != nil {
		return err
	}
	defer rpc_span.Finish()

	return nil
}

func f2(ptr_client *arpc.Client, ctx context.Context) error {
	rpc_span := rpc.New_rpc_span()
	if _, err := rpc_span.Start(ctx, "f2"); err != nil {
		return err
	}
	defer rpc_span.Finish()

	var reply string
	args := &Ping_args{Carrier: rpc_span.Carrier}

	if err := ptr_client.Call("v1.api.a_service.Ping", args, &reply); err != nil {
		return err
	}

	dev_log.Info("msg", reply)
	return nil
}

func f3(ctx context.Context) error {
	rpc_span := rpc.New_rpc_span()
	if _, err := rpc_span.Start(ctx, "f3"); err != nil {
		return err
	}
	defer rpc_span.Finish()
	return nil
}
