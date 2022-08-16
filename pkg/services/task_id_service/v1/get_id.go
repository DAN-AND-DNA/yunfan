package v1

import (
	"io"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_arpc "yunfan/sdk/arpc/task_id_service"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	Err_get_id_ok                 = pkg_errcode.New("get_id: ok", me, sdk_errcode.Code_s2s_ok)
	Err_get_id_need_tag           = pkg_errcode.New("get_id: need biz_tag", me, sdk_errcode.Code_s2s_need_arg)
	Err_get_id_disconnect_from_db = pkg_errcode.New("get_id: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (svc *Task_id_service) Get_id(rwc io.ReadWriteCloser, args *sdk_arpc.Get_id_args, reply *sdk_arpc.Get_id_reply) error {
	if args.Biz_tag == "" {
		reply.Err = Err_get_id_need_tag
		return nil
	}
	id, err := svc.srv.GetId(args.Biz_tag)
	if err != nil {
		reply.Err = Err_get_id_disconnect_from_db
		return nil
	}
	reply.Id = id
	reply.Err = Err_get_id_ok
	return nil
}
