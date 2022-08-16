package service

import (
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/services/task_id_service/error_code"
	sdk_errcode "yunfan/sdk/errcode"
)

type Service struct {
	alloc *Alloc
}

func NewService() (s *Service) {
	var err error
	s = &Service{}
	s.alloc, err = s.NewAllocId()
	if err != nil {

		log_err := pkg_errcode.New("new_alloc_id: Create "+err.Error(), error_code.Me, sdk_errcode.Code_db_internal_error)
		error_code.Print(log_err)
		panic(log_err)
	}

	return s
}
