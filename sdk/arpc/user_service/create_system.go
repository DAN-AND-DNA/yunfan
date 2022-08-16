package user_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Create_system_args struct {
	Tid             string `json:"tid" example:"事务id"`
	System_name     string `json:"system_name" example:"系统名称"`
	System_describe string `json:"system_describe" example:"系统描述"`
}

type Create_system_reply struct {
	Err *errcode.Errcode `json:"err"`
}

func Create_system(rpc_client *arpc.Client, args *Create_system_args, reply *Create_system_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return nil
	}
	svc_method := "v1.user-service.Create_system"

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
