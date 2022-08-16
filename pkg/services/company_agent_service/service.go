package main

import (
	pkg_rpc "yunfan/pkg/rpc"
	pkg_rpc_swag "yunfan/pkg/rpc/swag"
	v1 "yunfan/pkg/services/company_agent_service/v1"
)

var (
	service = v1.New()
)

func init() {
	pkg_rpc.Register_by_name(service, "v1.company-agent-service")
	pkg_rpc_swag.Register("POST", "/v1.company-agent-service.Get_company_data", service.Get_company_data_swag)
}
