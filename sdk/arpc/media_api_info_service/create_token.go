package media_api_info_service

import (
	"context"
	"time"
	pkg_errcode "yunfan/pkg/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Media_type_enum uint8

const (
	Media_empty Media_type_enum = iota
	Media_toutiao
)

type Create_token_args struct {
	Tid           string          `json:"tid" example:"事务id"`
	Media_type    Media_type_enum `json:"media_type" example:"1"` // 媒体类型 (1: 头条)
	Media_name    string          `json:"media_name" example:"媒体名"`
	Usage         string          `json:"usage" example:"用处"`
	App_id        uint64          `json:"appid" example:"1234567"` // 头条的appid
	Secret        string          `json:"secret" example:"头条密钥"`
	Refresh_token string          `json:"refresh_token" example:"头条刷新token"`
}

type Create_token_reply struct {
	Err *pkg_errcode.Errcode `json:"err"`
}

func Create_token(rpc_client *arpc.Client, args *Create_token_args, reply *Create_token_reply, timeout time.Duration) error {
	if rpc_client == nil {

		// this is mock
		return nil
	}
	svc_method := "v1.media-api-info-service.Create_token"

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
