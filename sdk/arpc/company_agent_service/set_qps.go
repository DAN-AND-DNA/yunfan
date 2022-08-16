package company_agent_service

import (
	pkg_errcode "yunfan/pkg/errcode"
)

type Set_qps_args struct {
	Qps             uint32 `json:"qps" example:"10"`    // 总qps
	Qps_api_toutiao uint32 `json:"qps_api" example:"2"` // 头条 api qps
}

type Set_qps_reply struct {
	Qps             uint32               `json:"qps" example:"10"`    // 总qps
	Qps_api_toutiao uint32               `json:"qps_api" example:"2"` // 头条 api qps
	Err             *pkg_errcode.Errcode `json:"err"`
}
