package log_service

import (
	"context"
	"time"

	"snk.git.node1/yunfan/arpc"
)

// Log_service.Insert_logs
func Insert_logs(rpc_client *arpc.Client, args *Insert_logs_args, reply *Insert_logs_reply, timeout time.Duration) error {
	if rpc_client == nil {
		return nil
	}

	if int64(timeout) == 0 {
		return rpc_client.Call("v1.api.log_service.Insert_logs", args, reply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case call := <-rpc_client.Go("v1.api.log_service.Insert_logs", args, reply, nil).Done:
		return call.Error
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil

}

type Insert_logs_args struct {
	Compress_type string // gzip
	Split_by      []byte
	Raw_json      []byte
}

type Insert_logs_reply struct {
}

// v1.api.log_service.Ping
func Ping(rpc_client *arpc.Client, args *Ping_args, reply *Ping_reply, timeout time.Duration) error {
	if rpc_client == nil {
		return nil
	}

	if int64(timeout) == 0 {
		return rpc_client.Call("v1.api.log_service.Ping", args, reply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case call := <-rpc_client.Go("v1.api.log_service.Ping", args, reply, nil).Done:
		return call.Error
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil

}

type Ping_args struct {
}

type Ping_reply struct {
	Message string
}
