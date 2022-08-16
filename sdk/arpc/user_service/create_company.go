package user_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Create_company_args struct {
	Tid              string `json:"tid" example:"事务id"`
	Sysid            uint64 `json:"sysid" example:"1"` // 系统id
	Company_name     string `json:"company_name" example:"公司名"`
	Company_describe string `json:"company_describe" example:"公司描述"`
}

type Create_company_reply struct {
	Err *errcode.Errcode `json:"err"`
}

func Create_company(rpc_client *arpc.Client, args *Create_company_args, reply *Create_company_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return nil
	}
	svc_method := "v1.user-service.Create_company"

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
