package main

import (
	pkg_dbs "yunfan/pkg/dbs"
	sdk_dbs "yunfan/sdk/dbs/task_id_service"
)

func init() {
	pkg_dbs.Register_postgres((&sdk_dbs.Segments{}).TableName(), "", &sdk_dbs.Segments{})
}
