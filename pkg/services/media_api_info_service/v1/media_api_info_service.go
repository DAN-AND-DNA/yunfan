package v1

import (
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_errcode "yunfan/sdk/errcode"

	"snk.git.node1/dan/go_request"
)

var (
	me = sdk_errcode.From_media_api_info_service
)

type Media_api_info_service struct {
	loc         *time.Location
	http_client go_request.Raw_request
}

func New() *Media_api_info_service {
	service := &Media_api_info_service{
		http_client: go_request.New(go_request.Config{
			Connect_timeout:                  3,
			Read_timeout:                     15,
			Write_timeout:                    15,
			Max_keepalive_idle_conn_duration: 15,
			Max_keepalive_conn_duration:      -1,
			Read_buffer_size:                 4096,
		}),
	}

	service.loc, _ = time.LoadLocation("Asia/Shanghai")

	return service
}

func Print(err *pkg_errcode.Errcode) {
	if err != nil {
		pkg_log.Error("to", me, "from", err.From(), "code", err.Code(), "msg", err.Error())
	}
}
