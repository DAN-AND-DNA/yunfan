package main

// 注册业务
import (
	pkg_rpc "yunfan/pkg/rpc"
	pkg_rpc_swag "yunfan/pkg/rpc/swag"
	v1 "yunfan/pkg/services/media_api_info_service/v1"
)

var (
	service = v1.New()
)

func init() {
	pkg_rpc.Register_by_name(service, "v1.media-api-info-service")
	pkg_rpc_swag.Register("POST", "/v1.media-api-info-service.Create_token", service.Create_token_swag)
	pkg_rpc_swag.Register("POST", "/v1.media-api-info-service.Manual_refresh_token", service.Manual_refresh_token_swag)
	pkg_rpc_swag.Register("POST", "/v1.media-api-info-service.Ping", service.Ping_swag)
}
