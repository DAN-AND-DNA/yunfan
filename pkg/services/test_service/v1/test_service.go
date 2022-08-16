package v1

import "time"

type Test_service struct {
	loc *time.Location
	qps uint32
}

func New() *Test_service {
	service := &Test_service{
		qps: 10,
	}

	service.loc, _ = time.LoadLocation("Asia/Shanghai")
	return service
}
