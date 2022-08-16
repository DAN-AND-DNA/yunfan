package v1

import (
	"context"
	"io"
	"strconv"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_arpc "yunfan/sdk/arpc/media_api_info_service"
	sdk_dbs "yunfan/sdk/dbs/media_api_info_service"
	sdk_errcode "yunfan/sdk/errcode"
	sdk_toutiao "yunfan/sdk/toutiao"

	"bytes"
	"net/http"

	"yunfan/pkg/test_tool"

	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	Err_manual_refresh_token_ok                 = pkg_errcode.New("manual_refresh_token: ok", me, sdk_errcode.Code_s2s_ok)
	Err_manual_refresh_token_need_tid           = pkg_errcode.New("manual_refresh_token: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_need_sid           = pkg_errcode.New("manual_refresh_token: need sub id", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_need_sname         = pkg_errcode.New("manual_refresh_token: need sub name", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_need_media_type    = pkg_errcode.New("manual_refresh_token: need media type", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_bad_media_type     = pkg_errcode.New("manual_refresh_token: bad media type", me, sdk_errcode.Code_s2s_bad_arg)
	Err_manual_refresh_token_need_media_name    = pkg_errcode.New("manual_refresh_token: need media name", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_bad_media_name     = pkg_errcode.New("manual_refresh_token: bad media name", me, sdk_errcode.Code_s2s_bad_arg)
	Err_manual_refresh_token_need_app_id        = pkg_errcode.New("manual_refresh_token: need app id", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_bad_app_id         = pkg_errcode.New("manual_refresh_token: bad app id", me, sdk_errcode.Code_s2s_bad_arg)
	Err_manual_refresh_token_need_secret        = pkg_errcode.New("manual_refresh_token: need secret", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_need_refresh_token = pkg_errcode.New("manual_refresh_token: need refresh token", me, sdk_errcode.Code_s2s_need_arg)
	Err_manual_refresh_token_disconnect_from_db = pkg_errcode.New("manual_refresh_token: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *Media_api_info_service) Manual_refresh_token(rwc io.ReadWriteCloser, args *sdk_arpc.Manual_refresh_token_args, reply *sdk_arpc.Manual_refresh_token_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_manual_refresh_token)
	sname := "manual_refresh_token"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_manual_refresh_token_need_tid
		return nil
	}

	switch args.Media_type {
	case sdk_arpc.Media_empty:
		reply.Err = Err_manual_refresh_token_need_media_type
		return nil
	case sdk_arpc.Media_toutiao:
		if args.App_id == 0 {
			reply.Err = Err_manual_refresh_token_need_app_id
			return nil
		}

	default:
		reply.Err = Err_manual_refresh_token_bad_media_type
		return nil
	}

	// 2. 验证事务和分支
	//db, ok := pkg_dbs.Postgres()
	db, ok := pkg_dbs.Begin_postgres()
	if !ok {
		reply.Err = Err_manual_refresh_token_disconnect_from_db
		Print(reply.Err)
		return nil
	}
	defer pkg_dbs.End_postgres()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid = ?", args.Tid).Find(&trans_records)
	cancel()

	if result.Error != nil {
		reply.Err = pkg_errcode.New("manual_refresh_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
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
			reply.Err = pkg_errcode.New("manual_refresh_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
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
		reply.Err = Err_manual_refresh_token_ok
		return nil
	}

	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	// 头条
	if args.Media_type == sdk_arpc.Media_toutiao {

		// 刷新token
		api_infos := []sdk_dbs.Toutiao_api_info{}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		result := db.WithContext(ctx).Where("appid = ?", args.App_id).Find(&api_infos)
		cancel()
		if result.Error != nil {
			reply.Err = pkg_errcode.New("manual_refresh_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}

		if len(api_infos) != 1 {
			reply.Err = Err_manual_refresh_token_bad_app_id
			Print(reply.Err)
			return nil
		}

		status_code, raw_resp, err := sdk_toutiao.Do_refresh_token(this.http_client, &sdk_toutiao.Refresh_token_req{
			App_id:        api_infos[0].App_id,
			Secret:        api_infos[0].Secret,
			Grant_type:    "refresh_token",
			Refresh_token: api_infos[0].Refresh_token}, 5)

		if err != nil {
			reply.Err = pkg_errcode.New("manual_refresh_token: "+err.Error(), me, sdk_errcode.Code_s2s_bad_arg)
			Print(reply.Err)
			return nil
		}

		if status_code > 299 {
			reply.Err = pkg_errcode.New("manual_refresh_token: status code is "+strconv.Itoa(status_code), me, sdk_errcode.Code_s2s_bad_arg)
			Print(reply.Err)
			return nil
		}

		if raw_resp.Code != 0 {
			reply.Err = pkg_errcode.New("manual_refresh_token: toutiao error msg is  "+raw_resp.Message, me, sdk_errcode.Code_s2s_bad_arg)
			Print(reply.Err)
			return nil
		}

		// 3. 一次性完成: 修改信息，确认分支，完成事务
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if result := tx.Create(&new_t_record); result.Error != nil {
				return result.Error
			}

			if result := tx.Model(&sdk_dbs.Toutiao_api_info{}).Where("app_id = ?", api_infos[0].App_id).Updates(sdk_dbs.Toutiao_api_info{Token_update_timestamp: (uint64)(now_time.Unix()), Token_expired_timestamp: raw_resp.Data.Expires_in + (uint64)(now_time.Unix()), Access_token: raw_resp.Data.Access_token, Refresh_token: raw_resp.Data.Refresh_token}); result.Error != nil {
				return result.Error
			}

			if result := tx.Model(&sdk_dbs.Transaction_info{}).Where("tid = ? AND sid = ? AND sname = ? ", args.Tid, "main", "main").Update("status", "done"); result.Error != nil {
				return result.Error
			}

			return nil
		})
		cancel()
		if err != nil {
			reply.Err = pkg_errcode.New("manual_refresh_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		}

		pkg_log.Info("manual_refresh", "toutiao token", "app_id", api_infos[0].App_id, "access_token", raw_resp.Data.Access_token, "refresh_token", raw_resp.Data.Refresh_token)
	}

	reply.Err = Err_manual_refresh_token_ok
	return nil
}

// @Summary 手动刷新token
// @Description 手动刷新token
// @Tags media-api-info-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Manual_refresh_token_args true "args"
// @Success 200 {object} sdk_arpc.Manual_refresh_token_reply "reply"
// @Router /v1.media-api-info-service.Manual_refresh_token [post]
func (this *Media_api_info_service) Manual_refresh_token_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Manual_refresh_token_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Manual_refresh_token_reply{}
	if err := this.Manual_refresh_token(&test_tool.Test_conn{}, args, reply); err != nil {
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
