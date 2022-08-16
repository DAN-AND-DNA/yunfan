package main

// 注册业务
import (
	pkg_rpc "yunfan/pkg/rpc"
	v1 "yunfan/pkg/services/task_id_service/v1"
)

var (
	service = v1.New()
)

func init() {
	pkg_rpc.Register_by_name(service, "v1.task-id-service")

}
