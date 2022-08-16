package auto_tasks

import (
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_errcode "yunfan/sdk/errcode"

	"snk.git.node1/dan/go_request"
)

var (
	me = sdk_errcode.From_auto_refresh_task
)

type Auto_tasks struct {
	http_client go_request.Raw_request
}

func New() *Auto_tasks {
	return &Auto_tasks{
		http_client: go_request.New(go_request.Config{
			Connect_timeout:                  3,
			Read_timeout:                     15,
			Write_timeout:                    15,
			Max_keepalive_idle_conn_duration: 15,
			Max_keepalive_conn_duration:      -1,
			Read_buffer_size:                 4096,
		}),
	}
}

func Print(err *pkg_errcode.Errcode) {
	if err != nil {
		pkg_log.Error("to", me, "from", err.From(), "code", err.Code(), "msg", err.Error())
	}
}
