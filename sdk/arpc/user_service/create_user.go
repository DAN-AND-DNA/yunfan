package user_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Create_user_args struct {
	Tid      string `json:"tid" example:"事务id"`
	Cid      uint64 `json:"cid" example:"1"` // 公司id
	Username string `json:"username" example:"用户名"`
	Password string `json:"password" example:"密码"`
}

type Create_user_reply struct {
	Err *errcode.Errcode `json:"err"`
}

func Create_user(rpc_client *arpc.Client, args *Create_user_args, reply *Create_user_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return nil
	}
	svc_method := "v1.user-service.Create_user"

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
