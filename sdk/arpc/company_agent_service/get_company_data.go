package company_agent_service

import (
	"context"
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_dbs "yunfan/sdk/dbs/company_agent_service"

	"snk.git.node1/yunfan/arpc"
)

type Get_company_data_args struct {
	Ids []uint64 `json:"ids"`
}

type Get_company_data_reply struct {
	Company_data_info_list []sdk_dbs.Company_data_info `json:"list"`
	Err                    *pkg_errcode.Errcode        `json:"err"`
}

func Get_company_data_raw(rpc_client *arpc.Client, args *Get_company_data_args, reply *Get_company_data_reply, timeout time.Duration) error {
	svc_method := "v1.company-agent-service.Get_company_data"

	if int64(timeout) == 0 {
		return rpc_client.Call(svc_method, args, reply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case call := <-rpc_client.Go(svc_method, args, reply, nil).Done:
		return call.Error
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func Get_company_data_mock(args *Get_company_data_args, reply *Get_company_data_reply) error {
	reply.Company_data_info_list = append(reply.Company_data_info_list, sdk_dbs.Company_data_info{
		Id:   1234567,
		Name: "测试公司2021",
		Cost: 120000,
	})
	return nil
}

func Get_company_data(rpc_client *arpc.Client, args *Get_company_data_args, reply *Get_company_data_reply, timeout time.Duration) error {
	if rpc_client == nil {
		// this is mock
		return Get_company_data_mock(args, reply)
	}

	return Get_company_data_raw(rpc_client, args, reply, timeout)
}
