package main

import (
	"os"
	"os/signal"
	"syscall"
	cmd "yunfan/cmd/media_api_info_service"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_log "yunfan/pkg/log"
	pkg_rpc "yunfan/pkg/rpc"
	pkg_swag "yunfan/pkg/rpc/swag"
	pkg_tasks "yunfan/pkg/tasks"

	_ "yunfan/pkg/services/media_api_info_service/docs"
)

// @title media-api-info-service
// @version 0.1.0
// @description 媒体api信息服务
// @termsOfService https://github.com/DAN-AND-DNA?tab=stars
// @contact.name Snk技术开发中心
// @contact.email danyang.chen@snkad.cn
// @BasePath /
func main() {
	var done = make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM)
	signal.Notify(done, syscall.SIGINT)

	// 1. get options from cmd
	cmd.Execute()
	env := cmd.Get_env()

	// 2. db
	if err := pkg_dbs.Open_postgres(env.DB_postgresql_dsn, false); err != nil {
		panic(err)
	}

	// 3. log
	if err := pkg_log.Connect(env.Log_level); err != nil {
		panic(err)
	}
	defer pkg_log.Shutdown(3)

	// 4. tasks
	pkg_tasks.Handle_task()
	defer pkg_tasks.Shutdown(3)

	// 5. rpc
	if err := pkg_rpc.Listen(env.Arpc_network_protocol, env.Arpc_app_protocol, env.Arpc_listen_port, env.Arpc_udp_key, env.Arpc_udp_salt, env.Arpc_conn_idle); err != nil {
		panic(err)
	}
	defer pkg_rpc.Shutdown(5)

	// 6. rpc swag
	if env.Arpc_swag == "on" {
		pkg_swag.Listen("v1.media-api-info-service", "3777")
	}

	_ = <-done
}
