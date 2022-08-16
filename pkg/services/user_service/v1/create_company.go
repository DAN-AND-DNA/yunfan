package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
	"yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/test_tool"
	sdk_arpc "yunfan/sdk/arpc/user_service"
	sdk_dbs "yunfan/sdk/dbs/user_service"
	sdk_errcode "yunfan/sdk/errcode"

	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	Err_create_company_ok                 = pkg_errcode.New("create_company: ok", me, sdk_errcode.Code_s2s_ok)
	Err_create_company_need_tid           = pkg_errcode.New("create_company: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_company_need_company_name  = pkg_errcode.New("create_company: need company name", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_company_need_sysid         = pkg_errcode.New("create_company: need system id", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_company_disconnect_from_db = pkg_errcode.New("create_company: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *User_service) Create_company(rwc io.ReadWriteCloser, args *sdk_arpc.Create_company_args, reply *sdk_arpc.Create_company_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_create_company)
	sname := "create_company"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_create_company_need_tid
		return nil
	}

	if args.Sysid == 0 {
		reply.Err = Err_create_company_need_sysid
		return nil
	}

	// 2. 验证事务和分支
	db, ok := dbs.Postgres()
	if !ok {
		reply.Err = Err_create_company_disconnect_from_db
		Print(reply.Err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid = ?", args.Tid).Find(&trans_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("create_company: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	has_registerd := false
	has_trans_done := false
	has_sub_done := false
	new_cid := uint64(0)

	if len(trans_records) > 0 {
		has_registerd = true
	}

	now_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), this.loc)
	if !has_registerd {
		// 未注册则注册事务
		new_record := &sdk_dbs.Transaction_info{
			Tid:              args.Tid,
			Sid:              "main",
			Sname:            "main",
			Status:           "registerd",
			Create_timestamp: (uint64)(now_time.Unix()),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if result := db.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_record); result.Error != nil {
			cancel()
			reply.Err = pkg_errcode.New("create_company: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}
		cancel()

	}

	for _, record := range trans_records {
		if record.Sid == "main" && record.Sname == "main" && record.Status == "done" {
			has_trans_done = true
			break

		}

		if record.Sid == sid {
			has_sub_done = true
		}

	}

	if has_trans_done || has_sub_done {
		reply.Err = Err_create_company_ok
		return nil
	}

	// 3. 获得公司id
	gen_id_args := sdk_arpc.Gen_id_args{
		Tid:  args.Tid,
		Type: sdk_arpc.Gen_company,
	}
	gen_id_reply := sdk_arpc.Gen_id_reply{}
	err := this.Gen_id(rwc, &gen_id_args, &gen_id_reply)
	if err != nil {
		reply.Err = pkg_errcode.New("create_company: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	if gen_id_reply.Err != nil {
		reply.Err = gen_id_reply.Err
		Print(reply.Err)
		return nil
	}
	new_cid = gen_id_reply.Id

	// 4. 一次完成: 创建公司，确认分支，完成事务
	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	new_c_record := &sdk_dbs.Company_info{
		Tid:              args.Tid,
		Cid:              new_cid,
		Company_name:     args.Company_name,
		Company_describe: args.Company_describe,
		Create_time:      now_time,
	}

	new_sc_record := &sdk_dbs.System_company_map{
		Sysid: args.Sysid,
		Cid:   new_cid,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(new_t_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_c_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_sc_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Model(&sdk_dbs.Transaction_info{}).Where("tid = ? AND sid = ? AND sname = ? ", args.Tid, "main", "main").Update("status", "done"); result.Error != nil {
			return result.Error
		}

		return nil

	})
	cancel()
	if err != nil {
		reply.Err = pkg_errcode.New("create_company: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	reply.Err = Err_create_company_ok
	return nil
}

// @Summary 创建公司
// @Description 创建公司
// @Tags user-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Create_company_args true "args"
// @Success 200 {object} sdk_arpc.Create_company_reply "reply"
// @Router /v1.user-service.Create_company [post]
func (this *User_service) Create_company_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Create_company_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Create_company_reply{}
	if err := this.Create_company(&test_tool.Test_conn{}, args, reply); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := enc.Encode(reply); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buffer.Bytes())
}
