package v1

import (
	"sync"
	"time"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	me = sdk_errcode.From_company_agent_service
)

type Company_agent_service struct {
	qps               uint32
	qps_api_toutiao   uint32
	loc               *time.Location
	transaction_infos sync.Map
}

func New() *Company_agent_service {
	service := &Company_agent_service{
		qps:             10,
		qps_api_toutiao: 2,
	}

	service.loc, _ = time.LoadLocation("Asia/Shanghai")
	return service
}
