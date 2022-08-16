package v1

import (
	"context"
	"io"
	"strconv"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	sdk_arpc "yunfan/sdk/arpc/media_api_info_service"
	sdk_dbs "yunfan/sdk/dbs/media_api_info_service"
	sdk_errcode "yunfan/sdk/errcode"

	"bytes"
	"net/http"
	"yunfan/pkg/test_tool"

	json "github.com/json-iterator/go"
	//"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	Err_create_token_ok                 = pkg_errcode.New("create_token: ok", me, sdk_errcode.Code_s2s_ok)
	Err_create_token_need_tid           = pkg_errcode.New("create_token: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_need_sid           = pkg_errcode.New("create_token: need sub id", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_need_sname         = pkg_errcode.New("create_token: need sub name", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_need_media_type    = pkg_errcode.New("create_token: need media type", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_bad_media_type     = pkg_errcode.New("create_token: bad media type", me, sdk_errcode.Code_s2s_bad_arg)
	Err_create_token_need_media_name    = pkg_errcode.New("create_token: need media name", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_bad_media_name     = pkg_errcode.New("create_token: bad media name", me, sdk_errcode.Code_s2s_bad_arg)
	Err_create_token_need_app_id        = pkg_errcode.New("create_token: need app id", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_need_secret        = pkg_errcode.New("create_token: need secret", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_need_refresh_token = pkg_errcode.New("create_token: need refresh token", me, sdk_errcode.Code_s2s_need_arg)
	Err_create_token_disconnect_from_db = pkg_errcode.New("create_token: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *Media_api_info_service) Create_token(rwc io.ReadWriteCloser, args *sdk_arpc.Create_token_args, reply *sdk_arpc.Create_token_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_create_token)
	sname := "create_token"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_create_token_need_tid
		return nil
	}

	switch args.Media_type {
	case sdk_arpc.Media_empty:
		reply.Err = Err_create_token_need_media_type
		return nil
	case sdk_arpc.Media_toutiao:
		if args.Media_name == "" {
			reply.Err = Err_create_token_need_media_name
			return nil
		}

		if args.Media_name != "今日头条" {
			reply.Err = Err_create_token_bad_media_name
			return nil
		}

		if args.App_id == 0 {
			reply.Err = Err_create_token_need_app_id
			return nil
		}

		if args.Secret == "" {
			reply.Err = Err_create_token_need_secret
			return nil
		}

		if args.Refresh_token == "" {
			reply.Err = Err_create_token_need_refresh_token
			return nil
		}

	default:
		reply.Err = Err_create_token_bad_media_type
		return nil
	}

	// 2. 验证事务和分支
	//db, ok := pkg_dbs.Postgres()
	db, ok := pkg_dbs.Begin_postgres()
	if !ok {
		reply.Err = Err_create_token_disconnect_from_db
		Print(reply.Err)
		return nil
	}
	defer pkg_dbs.End_postgres()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid=?", args.Tid).Find(&trans_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("create_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	has_registerd := false
	has_trans_done := false
	has_sub_done := false

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
			reply.Err = pkg_errcode.New("create_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
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
		reply.Err = Err_create_token_ok
		return nil
	}

	// 3. 一次性完成: 创建token信息表，确认分支，完成事务
	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	if args.Media_type == sdk_arpc.Media_toutiao {
		new_d_record := &sdk_dbs.Toutiao_api_info{
			Tid:                     args.Tid,
			App_id:                  args.App_id,
			Secret:                  args.Secret,
			Create_time:             now_time,
			Token_update_timestamp:  0,
			Token_expired_timestamp: 0,
			Refresh_token:           args.Refresh_token,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tx := db.WithContext(ctx).Begin()

		if err := tx.Create(new_t_record).Error; err != nil {
			tx.Rollback()
			reply.Err = pkg_errcode.New("create_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}

		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(new_d_record).Error; err != nil {
			tx.Rollback()
			reply.Err = pkg_errcode.New("create_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}

		if err := tx.Model(&sdk_dbs.Transaction_info{}).Where("tid = ? AND sid = ? AND sname = ?", args.Tid, "main", "main").Update("status", "done").Error; err != nil {
			tx.Rollback()
			reply.Err = pkg_errcode.New("create_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}

		if err := tx.Commit().Error; err != nil {
			reply.Err = pkg_errcode.New("create_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}
	}

	reply.Err = Err_create_token_ok
	return nil
}

// @Summary 创建token
// @Description 创建token
// @Tags media-api-info-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Create_token_args true "args"
// @Success 200 {object} sdk_arpc.Create_token_reply "reply"
// @Router /v1.media-api-info-service.Create_token [post]
func (this *Media_api_info_service) Create_token_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Create_token_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Create_token_reply{}
	if err := this.Create_token(&test_tool.Test_conn{}, args, reply); err != nil {
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
