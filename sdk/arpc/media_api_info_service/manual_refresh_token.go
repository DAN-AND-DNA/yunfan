package media_api_info_service

import (
	"context"
	"time"
	pkg_errcode "yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Manual_refresh_token_args struct {
	Tid        string          `json:"tid" example:"事务id"`
	Media_type Media_type_enum `json:"media_type" example:"1"`   //媒体类型 (1: 头条)
	App_id     uint64          `json:"app_id" example:"1234567"` // 头条app id
}

type Manual_refresh_token_reply struct {
	Err *pkg_errcode.Errcode `json:"err"`
}

func Manual_refresh_token(rpc_client *arpc.Client, args *Manual_refresh_token_args, reply *Manual_refresh_token_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return nil
	}
	svc_method := "v1.media-api-info-service.Manual_refresh_token"

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
