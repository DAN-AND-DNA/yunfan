package main

import (
	pkg_http "yunfan/pkg/http"
	v1 "yunfan/pkg/services/test_service/v1"
)

var (
	service = v1.New()
)

func init() {
	pkg_http.Register("GET", "ping", service.Ping)
}
