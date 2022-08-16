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
	Err_create_user_ok                 = pkg_errcode.New("create_user: ok", me, sdk_errcode.Code_s2s_ok)
	Err_create_user_need_tid           = pkg_errcode.New("create_user: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_user_need_username      = pkg_errcode.New("create_user: need username", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_user_need_cid           = pkg_errcode.New("create_user: need company id", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_user_disconnect_from_db = pkg_errcode.New("create_user: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *User_service) Create_user(rwc io.ReadWriteCloser, args *sdk_arpc.Create_user_args, reply *sdk_arpc.Create_user_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_create_user)
	sname := "create_user"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_create_user_need_tid
		return nil
	}

	if args.Cid == 0 {
		reply.Err = Err_create_user_need_cid
		return nil
	}

	// 2. 验证事务和分支
	db, ok := dbs.Postgres()
	if !ok {
		reply.Err = Err_create_user_disconnect_from_db
		Print(reply.Err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid = ?", args.Tid).Find(&trans_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("create_user: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	has_registerd := false
	has_trans_done := false
	has_sub_done := false
	new_uid := uint64(0)

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

		result := db.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_record)

		cancel()
		if result.Error != nil {
			reply.Err = pkg_errcode.New("create_user: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}
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
		reply.Err = Err_create_user_ok
		return nil
	}

	// 3. 获得用户id
	gen_id_args := sdk_arpc.Gen_id_args{
		Tid:  args.Tid,
		Type: sdk_arpc.Gen_user,
	}
	gen_id_reply := sdk_arpc.Gen_id_reply{}
	err := this.Gen_id(rwc, &gen_id_args, &gen_id_reply)
	if err != nil {
		reply.Err = pkg_errcode.New("create_user: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}
	if gen_id_reply.Err != nil {
		reply.Err = gen_id_reply.Err
		Print(reply.Err)
		return nil
	}
	new_uid = gen_id_reply.Id

	// 4. 一次完成: 创建用户，确认分支，完成事务
	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	new_u_record := &sdk_dbs.User_info{
		Tid:         args.Tid,
		Uid:         new_uid,
		Username:    args.Username,
		Create_time: now_time,
	}

	new_cu_record := &sdk_dbs.Company_user_map{
		Cid: args.Cid,
		Uid: new_uid,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if result := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_u_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_cu_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Create(new_t_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Model(&sdk_dbs.Transaction_info{}).Where("tid = ? AND sid = ? AND sname = ? ", args.Tid, "main", "main").Update("status", "done"); result.Error != nil {
			return result.Error
		}

		return nil

	})
	cancel()
	if err != nil {
		reply.Err = pkg_errcode.New("create_user: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	reply.Err = Err_create_user_ok
	return nil
}

// @Summary 创建用户
// @Description 创建用户
// @Tags user-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Create_user_args true "args"
// @Success 200 {object} sdk_arpc.Create_user_reply "reply"
// @Router /v1.user-service.Create_user [post]
func (this *User_service) Create_user_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Create_user_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Create_user_reply{}
	if err := this.Create_user(&test_tool.Test_conn{}, args, reply); err != nil {
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
