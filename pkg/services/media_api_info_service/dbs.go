package main

import (
	pkg_dbs "yunfan/pkg/dbs"
	sdk_dbs "yunfan/sdk/dbs/media_api_info_service"
)

func init() {
	pkg_dbs.Register_postgres("transaction_info", "WITH (fillfactor=50)", &sdk_dbs.Transaction_info{}) // 事务表
	pkg_dbs.Register_postgres("toutiao_api_info", "WITH (fillfactor=70)", &sdk_dbs.Toutiao_api_info{}) // 头条信息表
}
