package main

import (
	"os"
	"os/signal"
	"syscall"
	cmd "yunfan/cmd/task_id_service"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_log "yunfan/pkg/log"
	pkg_rpc "yunfan/pkg/rpc"
)

// @title task-id-service
// @version 0.1.0
// @description 媒体api信息服务
// @contact.name Snk技术开发中心
// @contact.email youhong.yang@snkad.cn
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

	// 4. rpc
	if err := pkg_rpc.Listen(env.Arpc_network_protocol, env.Arpc_app_protocol, env.Arpc_listen_port, env.Arpc_udp_key, env.Arpc_udp_salt, env.Arpc_conn_idle); err != nil {
		panic(err)
	}
	defer pkg_rpc.Shutdown(5)

	//5.etcd
	if err := Init_master_node(env.Etcd_addrs, env.Arpc_listen_port, 30); err != nil {
		panic(err)
	}
	defer Master_node.Close_apply_master_node()

	service.Init_service()

	<-done
}
