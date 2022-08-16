package main

import (
	pkg_dbs "yunfan/pkg/dbs"
	sdk_dbs "yunfan/sdk/dbs/user_service"
)

func init() {
	pkg_dbs.Register_postgres("transaction_info", "WITH(fillfactor=50)", &sdk_dbs.Transaction_info{})
	pkg_dbs.Register_postgres("system_info", "WITH(fillfactor=50)", &sdk_dbs.System_info{})
	pkg_dbs.Register_postgres("system_company_map", "WITH(fillfactor=50)", &sdk_dbs.System_company_map{})
	pkg_dbs.Register_postgres("company_info_id", "WITH(fillfactor=50)", &sdk_dbs.Company_info_id{})
	pkg_dbs.Register_postgres("company_info", "WITH(fillfactor=50)", &sdk_dbs.Company_info{})
	pkg_dbs.Register_postgres("company_user_map", "WITH(fillfactor=50)", &sdk_dbs.Company_user_map{})
	pkg_dbs.Register_postgres("user_info_id", "WITH(fillfactor=50)", &sdk_dbs.User_info_id{})
	pkg_dbs.Register_postgres("user_info", "WITH(fillfactor=50)", &sdk_dbs.User_info{})
}
