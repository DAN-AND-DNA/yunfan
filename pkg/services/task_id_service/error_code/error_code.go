package error_code

import (
	pkg_errcode "yunfan/pkg/errcode"
	pkg_log "yunfan/pkg/log"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	Err_get_all_segments_disconnect_from_db = pkg_errcode.New("get_all_segement: disconnect from db", Me, sdk_errcode.Code_db_internal_error)
	Me                                      = sdk_errcode.From_media_api_info_service
)

func Print(err *pkg_errcode.Errcode) {
	if err != nil {
		pkg_log.Error("to", Me, "from", err.From(), "code", err.Code(), "msg", err.Error())
	}
}
