package v1

import (
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	"yunfan/pkg/services/task_id_service/v1/service"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	me = sdk_errcode.From_task_id_service
)

type Task_id_service struct {
	loc *time.Location
	srv *service.Service
}

func New() *Task_id_service {
	svc := &Task_id_service{}
	svc.loc, _ = time.LoadLocation("Asia/Shanghai")
	return svc
}

func (svc *Task_id_service) Init_service() {
	svc.srv = service.NewService()
}

func Print(err *pkg_errcode.Errcode) {
	if err != nil {
		pkg_log.Error("to", me, "from", err.From(), "code", err.Code(), "msg", err.Error())
	}
}
