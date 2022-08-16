package auto_tasks

import (
	"context"
	"strconv"
	"time"
	"yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_dbs "yunfan/sdk/dbs/media_api_info_service"
	sdk_errcode "yunfan/sdk/errcode"
	sdk_toutiao "yunfan/sdk/toutiao"
)

var (
	Err_refresh_toutiao_token_disconnect_from_db = pkg_errcode.New("refresh_toutiao_token: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *Auto_tasks) Refresh_toutiao_token(current_execute_num uint64) error {
	// get info from db
	api_infos := []sdk_dbs.Toutiao_api_info{}

	db, ok := dbs.Postgres()
	if !ok {
		Print(Err_refresh_toutiao_token_disconnect_from_db)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	result := db.WithContext(ctx).Find(&api_infos)
	cancel()
	if result.Error != nil {
		Print(pkg_errcode.New("refresh_toutiao_token: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error))
		return nil
	}

	if result.RowsAffected == 0 {
		return nil
	}

	for _, api_info := range api_infos {
		if api_info.Secret == "" || api_info.Refresh_token == "" {
			Print(pkg_errcode.New("refresh_toutiao_token: "+strconv.FormatUint(api_info.App_id, 10)+" no secret or no refresh token", me, sdk_errcode.Code_s2s_bad_arg))
			continue
		}

		now_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)

		if (uint64)(now_time.Unix()) >= api_info.Token_expired_timestamp {
			//refresh
			raw_req := sdk_toutiao.Refresh_token_req{
				App_id:        api_info.App_id,
				Secret:        api_info.Secret,
				Grant_type:    "refresh_token",
				Refresh_token: api_info.Refresh_token,
			}

			status_code, raw_resp, err := sdk_toutiao.Do_refresh_token(this.http_client, &raw_req, 7)
			if err != nil {
				Print(pkg_errcode.New("refresh_toutiao_token: "+err.Error(), me, sdk_errcode.Code_s2s_bad_arg))
				continue
			}

			if status_code > 299 {
				Print(pkg_errcode.New("refresh_toutiao_token: status code is "+strconv.Itoa(status_code), me, sdk_errcode.Code_s2s_bad_arg))
				continue
			}

			if raw_resp.Code != 0 {
				Print(pkg_errcode.New("refresh_toutiao_token: toutiao error msg is "+raw_resp.Message, me, sdk_errcode.Code_s2s_bad_arg))
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			result := db.WithContext(ctx).Model(&sdk_dbs.Toutiao_api_info{}).Where("app_id = ?", api_info.App_id).Updates(sdk_dbs.Toutiao_api_info{Token_update_timestamp: (uint64)(now_time.Unix()), Token_expired_timestamp: raw_resp.Data.Expires_in + (uint64)(now_time.Unix()), Access_token: raw_resp.Data.Access_token, Refresh_token: raw_resp.Data.Refresh_token})
			cancel()

			if result.Error != nil {
				Print(pkg_errcode.New("refresh_toutiao_token: "+err.Error(), me, sdk_errcode.Code_db_internal_error))
				continue
			}

			pkg_log.Info("auto_refresh", "toutiao token", "app_id", api_info.App_id, "access_token", raw_resp.Data.Access_token, "refresh_token", raw_resp.Data.Refresh_token)

		}
	}

	return nil
}
