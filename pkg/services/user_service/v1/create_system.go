package v1

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
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
	Err_create_system_ok                 = pkg_errcode.New("create_system: ok", me, sdk_errcode.Code_s2s_ok)
	Err_create_system_need_tid           = pkg_errcode.New("create_system: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_system_need_system_name   = pkg_errcode.New("create_system: need system name", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_system_disconnect_from_db = pkg_errcode.New("create_system: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *User_service) Create_system(rwc io.ReadWriteCloser, args *sdk_arpc.Create_system_args, reply *sdk_arpc.Create_system_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_create_system)
	sname := "create_system"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_create_system_need_tid
		return nil
	}

	if args.System_name == "" {
		reply.Err = Err_create_system_need_system_name
		return nil
	}

	// 2. 验证事务和分支
	db, ok := pkg_dbs.Postgres()
	if !ok {
		reply.Err = Err_create_system_disconnect_from_db
		Print(reply.Err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid = ?", args.Tid).Find(&trans_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("create_system: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	has_registerd := false
	has_trans_done := false
	has_sub_done := false

	if len(trans_records) > 0 {
		log.Println(len(trans_records))
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
			reply.Err = pkg_errcode.New("create_system: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
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
		reply.Err = Err_create_system_ok
		return nil
	}

	// 3. 一次性完成: 创建系统，确认分支，完成事务
	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	new_s_record := &sdk_dbs.System_info{
		Tid:             args.Tid,
		System_name:     args.System_name,
		System_describe: args.System_describe,
		Create_time:     now_time,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result = tx.Create(&new_t_record)
		if result.Error != nil {
			return result.Error
		}

		if result := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_s_record); result.Error != nil {
			return result.Error
		}

		if result := tx.Model(&sdk_dbs.Transaction_info{}).Where("tid = ? AND sid = ? AND sname = ?", args.Tid, "main", "main").Update("status", "done"); result.Error != nil {
			return result.Error
		}

		return nil

	})
	cancel()
	if err != nil {
		reply.Err = pkg_errcode.New("create_system: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	reply.Err = Err_create_system_ok
	return nil
}

// @Summary 创建系统
// @Description 创建系统
// @Tags user-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Create_system_args true "args"
// @Success 200 {object} sdk_arpc.Create_system_reply "reply"
// @Router /v1.user-service.Create_system [post]
func (this *User_service) Create_system_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Create_system_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Create_system_reply{}
	if err := this.Create_system(&test_tool.Test_conn{}, args, reply); err != nil {
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
