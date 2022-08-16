package v1

import (
	"context"
	"errors"
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_rpc_media_api_info_service "yunfan/sdk/arpc/media_api_info_service"
	sdk_errcode "yunfan/sdk/errcode"

	"net"
	pkg_rpc "yunfan/pkg/rpc"

	"snk.git.node1/yunfan/arpc"
)

var (
	me = sdk_errcode.From_user_service
)

type User_service struct {
	loc                         *time.Location
	ok_media_api_info_service   bool
	pool_media_api_info_service pkg_rpc.Std_rpc_client_pool
}

func New(raw_is_mock ...bool) *User_service {
	service := &User_service{}
	service.loc, _ = time.LoadLocation("Asia/Shanghai")
	var err error

	is_mock := false
	if len(raw_is_mock) != 0 {
		is_mock = raw_is_mock[0]
	}

	if is_mock {
		service.pool_media_api_info_service, err = pkg_rpc.New_json_rpc_client_pool(2, 5,

			// 1. conn_factory
			func() (net.Conn, error) { return nil, nil },

			// 2. keep_alive_callback
			func(rpc_client *arpc.Client) bool {
				return true
			}, 15)
		if err != nil {
			panic(err)
		}

		return service
	} else {
		service.pool_media_api_info_service, err = pkg_rpc.New_json_rpc_client_pool(2, 5,

			// 1. conn_factory
			func() (net.Conn, error) {
				tcp_conn, err := net.DialTimeout("tcp", "media-api-info-service:37001", 3*time.Second)
				if err != nil {
					return nil, err
				}
				return tcp_conn, nil
			},

			// 2. keep_alive_callback
			func(rpc_client *arpc.Client) bool {
				for i := 0; i < 2; i++ {
					args := sdk_rpc_media_api_info_service.Ping_args{Msg: "ping"}
					reply := sdk_rpc_media_api_info_service.Ping_reply{}
					err = sdk_rpc_media_api_info_service.Ping(rpc_client, &args, &reply, 5*time.Second)
					if err == nil {
						return true
					}

					if errors.Is(err, context.DeadlineExceeded) {
						continue
					} else {
						return false
					}
				}
				return false
			}, 15)
		if err != nil {
			panic(err)
		}
	}

	return service
}

func Print(err *pkg_errcode.Errcode) {
	if err != nil {
		pkg_log.Error("to", me, "from", err.From(), "code", err.Code(), "msg", err.Error())
	}
}
