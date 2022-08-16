package user_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Gen_type_enum uint8

const (
	Gen_empty Gen_type_enum = iota
	Gen_company
	Gen_user
)

type Gen_id_args struct {
	Tid  string        `json:"tid" example:"事务id"`
	Type Gen_type_enum `json:"type" example:"1"` // 1: 公司 2: 用户
}

type Gen_id_reply struct {
	Id  uint64 `json:"id" example:"1"` // 产生的id
	Err *errcode.Errcode
}

func Gen_id(rpc_client *arpc.Client, args *Gen_id_args, reply *Gen_id_reply, timeout time.Duration) error {
	if rpc_client == nil {
		//this is mock
		return nil
	}
	svc_method := "v1.user-service.Gen_id"

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
