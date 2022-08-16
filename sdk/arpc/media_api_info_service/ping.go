package media_api_info_service

import (
	"context"
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_errcode "yunfan/sdk/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Ping_args struct {
	Msg string `json:"msg" example:"ping"`
}

type Ping_reply struct {
	Msg string               `json:"msg" example:"pong"`
	Err *pkg_errcode.Errcode `json:"err"`
}

func Ping_raw(rpc_client *arpc.Client, args *Ping_args, reply *Ping_reply, timeout time.Duration) error {
	svc_method := "v1.media-api-info-service.Ping"

	if int64(timeout) == 0 {
		return rpc_client.Call(svc_method, args, reply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case call := <-rpc_client.Go(svc_method, args, reply, nil).Done:
		return call.Error
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func Ping_mock(args *Ping_args, reply *Ping_reply) error {
	reply.Msg = "pong"
	reply.Err = pkg_errcode.New("ping: ok", sdk_errcode.From_media_api_info_service, sdk_errcode.Code_s2s_ok)
	return nil
}

func Ping(rpc_client *arpc.Client, args *Ping_args, reply *Ping_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return Ping_mock(args, reply)
	}

	// do request
	return Ping_raw(rpc_client, args, reply, timeout)
}
