package main

import (
	pkg_dbs "yunfan/pkg/dbs"
	sdk_dbs "yunfan/sdk/dbs/company_agent_service"
)

func init() {
	pkg_dbs.Register_sqlite("transaction_info", "", &sdk_dbs.Transaction_info{})
	pkg_dbs.Register_sqlite("company_data_info", "", &sdk_dbs.Company_data_info{})
}
