package v1

import (
	"io"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_arpc "yunfan/sdk/arpc/company_agent_service"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	Err_set_qps_ok = pkg_errcode.New("set_qps: ok", me, sdk_errcode.Code_s2s_ok)
)

func (this *Company_agent_service) Set_qps(rwc io.ReadWriteCloser, args *sdk_arpc.Set_qps_args, reply *sdk_arpc.Set_qps_reply) error {
	this.qps = args.Qps
	this.qps_api_toutiao = args.Qps_api_toutiao
	reply.Qps = this.qps
	reply.Qps_api_toutiao = this.qps_api_toutiao

	reply.Err = Err_set_qps_ok
	return nil
}
