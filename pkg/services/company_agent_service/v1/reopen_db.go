package v1

import (
	"io"
	"strconv"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_arpc "yunfan/sdk/arpc/company_agent_service"
	sdk_dbs "yunfan/sdk/dbs/company_agent_service"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	Err_reopen_db_ok       = pkg_errcode.New("reopen_db: ok", me, sdk_errcode.Code_s2s_ok)
	Err_reopen_db_retry    = pkg_errcode.New("reopen_db: retry", me, sdk_errcode.Code_s2s_retry)
	Err_reopen_db_need_tid = pkg_errcode.New("reopen_db: need tid", me, sdk_errcode.Code_s2s_need_arg)
)

func (this *Company_agent_service) Reopen_db(rwc io.ReadWriteCloser, args *sdk_arpc.Reopen_db_args, reply *sdk_arpc.Reopen_db_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_reopen_db)
	sname := "reopen_db"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_reopen_db_need_tid
		return nil
	}

	now_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), this.loc)

	if _, is_exist := this.transaction_infos.LoadOrStore(args.Tid+"_"+sid+"_registerd", &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Status:           "registerd",
		Create_timestamp: (uint64)(now_time.Unix()),
	}); is_exist {

		// 2. 已被注册获得结果
		if raw_val, ok := this.transaction_infos.Load(args.Tid + "_" + sid + "_done"); ok {
			trans_info := raw_val.(*sdk_dbs.Transaction_info)

			if trans_info.Status == "done" {
				reply.Err = Err_reopen_db_ok
				return nil
			} else {
				// error
				reply.Err = pkg_errcode.New("reopen_db: "+trans_info.Status, me, sdk_errcode.Code_db_internal_error)
				return nil

			}
		} else {
			reply.Err = Err_reopen_db_retry
			return nil
		}

	}

	// 3. 任务失败
	if err := pkg_dbs.Reopen_sqlite(); err != nil {
		this.transaction_infos.Store(args.Tid+"_"+sid+"_done", &sdk_dbs.Transaction_info{
			Tid:              args.Tid,
			Sid:              sid,
			Sname:            sname,
			Status:           err.Error(),
			Create_timestamp: (uint64)(now_time.Unix()),
		})

		reply.Err = pkg_errcode.New("reopen_db: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		return nil
	}

	// 4. 任务成功
	this.transaction_infos.Store(args.Tid+"_"+sid+"_done", &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Status:           "done",
		Create_timestamp: (uint64)(now_time.Unix()),
	})

	reply.Err = Err_reopen_db_ok
	return nil
}
