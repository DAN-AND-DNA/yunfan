package main

// 注册业务
import (
	"yunfan/pkg/app-example/services/v1/api/a_service"
	pkg_rpc "yunfan/pkg/rpc"
)

var (
	default_a_service = a_service.New_a_service()
)

func init() {
	// 注册

	// rpc api
	pkg_rpc.Register_by_name(default_a_service, "v1.api.a_service")
}
