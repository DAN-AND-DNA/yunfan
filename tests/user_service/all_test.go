package user_service

import (
	"fmt"
	"net"
	"testing"
	"time"

	pkg_dbs "yunfan/pkg/dbs"
	pkg_rpc "yunfan/pkg/rpc"
	v1 "yunfan/pkg/services/user_service/v1"
	sdk_rpc "yunfan/sdk/arpc/user_service"
	sdk_errcode "yunfan/sdk/errcode"

	"github.com/stretchr/testify/require"
	"snk.git.node1/dan/yoursql"
	"snk.git.node1/yunfan/arpc/jsonrpc"
)

var (
	c_conn, s_conn = net.Pipe()
	rpc_client     = jsonrpc.NewClient(c_conn)
)

func init() {
	// 1. mock db
	if err := pkg_dbs.Open_postgres("dsn", true); err != nil {
		panic(err)
	}

	// 2. init target rpc server
	service := v1.New(true)
	rpc_server := pkg_rpc.New_server()
	if err := rpc_server.Register_by_name(service, "v1.user-service"); err != nil {
		panic(err)
	}
	go rpc_server.Test_json(s_conn)
}

func Test_user_service_create_system(t *testing.T) {
	r := require.New(t)

	// 1. mock db result
	sql_key1 := `SELECT * FROM "transaction_infos" WHERE tid = 事务id`
	yoursql.Set_expect(0, 0, sql_key1, []map[string]interface{}{
		{},
	})
	defer yoursql.Clean_expect(sql_key1)

	now_time, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	r.Nil(err)

	sql_key2 := fmt.Sprintf(`INSERT INTO "transaction_infos" ("tid","sid","sname","status","create_timestamp") VALUES (事务id,main,main,registerd,%d) ON CONFLICT DO NOTHING RETURNING "id"`, now_time.Unix())

	yoursql.Set_expect(0, 0, sql_key2, []map[string]interface{}{
		{"id": 7},
	})

	defer yoursql.Clean_expect(sql_key2)

	sql_key3 := fmt.Sprintf(`INSERT INTO "transaction_infos" ("tid","sid","sname","status","create_timestamp") VALUES (事务id,3001,create_system,,%d) RETURNING "id"`, now_time.Unix())
	yoursql.Set_expect(0, 0, sql_key3, []map[string]interface{}{
		{"id": 7},
	})
	defer yoursql.Clean_expect(sql_key3)

	sql_key4 := fmt.Sprintf(`INSERT INTO "system_infos" ("tid","system_name","system_describe","create_time") VALUES (事务id,mkt,商业版,%s +0800 CST) ON CONFLICT DO NOTHING RETURNING "sysid"`, now_time.Format("2006-01-02 15:04:05"))
	yoursql.Set_expect(0, 0, sql_key4, []map[string]interface{}{
		{"sysid": 37},
	})
	defer yoursql.Clean_expect(sql_key4)

	sql_key5 := `UPDATE "transaction_infos" SET "status"=done WHERE tid = 事务id AND sid = main AND sname = main`
	yoursql.Set_expect(0, 0, sql_key5, []map[string]interface{}{
		{},
	})
	defer yoursql.Clean_expect(sql_key5)

	// 2. do rpc request
	create_system_args := sdk_rpc.Create_system_args{
		Tid:             "事务id",
		System_name:     "mkt",
		System_describe: "商业版",
	}
	create_system_reply := sdk_rpc.Create_system_reply{}
	err = sdk_rpc.Create_system(rpc_client, &create_system_args, &create_system_reply, 0)
	r.Nil(err)
	r.Equal(create_system_reply.Err.Msg, v1.Err_create_system_ok.Msg)
	r.Equal(sdk_errcode.From_user_service, create_system_reply.Err.From())
	r.Equal(sdk_errcode.Code_s2s_ok, create_system_reply.Err.Code())
}
