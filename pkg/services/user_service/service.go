package main

// 注册业务
import (
	pkg_rpc "yunfan/pkg/rpc"
	pkg_rpc_swag "yunfan/pkg/rpc/swag"
	v1 "yunfan/pkg/services/user_service/v1"
)

var (
	service = v1.New()
)

func init() {
	pkg_rpc.Register_by_name(service, "v1.user-service")
	pkg_rpc_swag.Register("POST", "/v1.user-service.Create_system", service.Create_system_swag)
	pkg_rpc_swag.Register("POST", "/v1.user-service.Create_company", service.Create_company_swag)
	pkg_rpc_swag.Register("POST", "/v1.user-service.Create_user", service.Create_user_swag)
	pkg_rpc_swag.Register("POST", "/v1.user-service.Gen_id", service.Gen_id_swag)
	pkg_rpc_swag.Register("POST", "/v1.user-service.Get_auth_info", service.Get_auth_info_swag)
	pkg_rpc_swag.Register("POST", "/v1.user-service.Get_idle_toutiao_appid", service.Get_idle_toutiao_appid_swag)
}
