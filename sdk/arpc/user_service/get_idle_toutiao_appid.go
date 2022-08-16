package user_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Get_idle_toutiao_appid_args struct {
}

type Get_idle_toutiao_appid_reply struct {
	Err *errcode.Errcode
}

func Get_idle_toutiao_appid(rpc_client *arpc.Client, args *Get_idle_toutiao_appid_args, reply *Get_idle_toutiao_appid_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return nil
	}
	svc_method := "v1.user-service.Get_idle_toutiao_appid"

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
