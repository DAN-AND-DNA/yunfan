package company_agent_service

import (
	pkg_errcode "yunfan/pkg/errcode"
)

type Reopen_db_args struct {
	Tid string `json:"tid"`
}

type Reopen_db_reply struct {
	Err *pkg_errcode.Errcode `json:"err"`
}
