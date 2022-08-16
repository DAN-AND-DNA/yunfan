package main

import (
	"os"
	"os/signal"
	"syscall"
	cmd "yunfan/cmd/app-example"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_http "yunfan/pkg/http"
	pkg_log "yunfan/pkg/log"
	pkg_rpc "yunfan/pkg/rpc"
	pkg_tasks "yunfan/pkg/tasks"
)

func main() {
	var done = make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM)
	signal.Notify(done, syscall.SIGINT)

	// get options from cmd

	cmd.Execute()
	env := cmd.Get_env()

	if err := pkg_rpc.Set_remote_tracer(env.Tracer_service_name, env.Tracer_remote_hostname, env.Tracer_remote_port); err != nil {
		panic(err)
	}

	defer pkg_rpc.Shutdown_tracer()

	if err := pkg_log.Connect(env.Log_level, env.Log_remote_host, env.Log_remote_port, env.Log_remote_heartbeat_interval, env.Log_remote_flush_interval); err != nil {
		panic(err)
	}

	defer pkg_log.Shutdown(7)

	pkg_http_server_config := pkg_http.Http_server_config{
		Run_as:             env.Run_as,
		Port:               env.Http_listen_port,
		Allow_origins:      env.Http_allow_origins,
		Server_name:        env.Http_header_name,
		Body_limit_bytes:   env.Http_body_limit_bytes,
		Read_timeout_secs:  env.Http_read_timeout_secs,
		Write_timeout_secs: env.Http_write_timeout_secs,
		Max_conns:          env.Http_max_conns,
		Get_only:           true,
	}
	if err := pkg_http.Listen(pkg_http_server_config); err != nil {
		panic(err)
	}
	defer pkg_http.Shutdown(5)

	pkg_tasks.Handle_task()
	defer pkg_tasks.Shutdown(3)

	if err := pkg_dbs.Open_ck(env.DB_clickhouse_dns); err != nil {
		panic(err)
	}

	if err := pkg_rpc.Listen(env.Rpc_protocol_type, env.Rpc_app_protocol_type, env.Rpc_listen_port, env.Rpc_key, env.Rpc_salt, env.Rpc_idle); err != nil {
		panic(err)
	}
	defer pkg_rpc.Shutdown(5)

	_ = <-done
}
