package media_api_info_service

import (
	"net"
	"testing"
	sdk_rpc "yunfan/sdk/arpc/media_api_info_service"

	pkg_dbs "yunfan/pkg/dbs"
	pkg_rpc "yunfan/pkg/rpc"
	v1 "yunfan/pkg/services/media_api_info_service/v1"

	"github.com/stretchr/testify/require"

	//sdk_rpc_env "yunfan/sdk/arpc"
	sdk_errcode "yunfan/sdk/errcode"

	"fmt"
	"time"

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
	service := v1.New()
	rpc_server := pkg_rpc.New_server()
	if err := rpc_server.Register_by_name(service, "v1.media-api-info-service"); err != nil {
		panic(err)
	}
	go rpc_server.Test_json(s_conn)

}

func Test_media_api_info_service_Ping(t *testing.T) {
	r := require.New(t)

	ping_args := sdk_rpc.Ping_args{Msg: "ping"}
	ping_reply := sdk_rpc.Ping_reply{}
	err := sdk_rpc.Ping_raw(rpc_client, &ping_args, &ping_reply, 0)
	r.Nil(err)
	r.Equal(sdk_errcode.From_media_api_info_service, ping_reply.Err.From())
	r.Equal(sdk_errcode.Code_s2s_ok, ping_reply.Err.Code())
	r.Equal("pong", ping_reply.Msg)
}

func Test_media_api_info_service_Create_token(t *testing.T) {
	r := require.New(t)

	// 1. mock db result
	sql_key1 := `SELECT * FROM "transaction_infos" WHERE tid=事务id`
	yoursql.Set_expect(0, 0, sql_key1, []map[string]interface{}{
		{},
	})
	defer yoursql.Clean_expect(sql_key1)

	now_time, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	r.Nil(err)

	sql_key2 := fmt.Sprintf(`INSERT INTO "transaction_infos" ("tid","sid","sname","status","create_timestamp") VALUES (事务id,main,main,registerd,%d) ON CONFLICT DO NOTHING RETURNING "id"`, now_time.Unix())
	yoursql.Set_expect(0, 0, sql_key2, []map[string]interface{}{
		{"id": 37},
	})
	defer yoursql.Clean_expect(sql_key2)

	sql_key3 := fmt.Sprintf(`INSERT INTO "transaction_infos" ("tid","sid","sname","status","create_timestamp") VALUES (事务id,2001,create_token,,%d) RETURNING "id"`, now_time.Unix())
	yoursql.Set_expect(0, 0, sql_key3, []map[string]interface{}{
		{"id": 7},
	})
	defer yoursql.Clean_expect(sql_key3)

	sql_key4 := fmt.Sprintf(`INSERT INTO "toutiao_api_infos" ("tid","secret","create_time","token_update_timestamp","token_expired_timestamp","access_token","refresh_token","app_id") VALUES (事务id,this is a secret,%s +0800 CST,0,0,,xxxrftoken,1234567) ON CONFLICT DO NOTHING RETURNING "app_id"`, now_time.Format("2006-01-02 15:04:05"))
	yoursql.Set_expect(0, 0, sql_key4, []map[string]interface{}{
		{"app_id": 1234567},
	})
	defer yoursql.Clean_expect(sql_key4)

	sql_key5 := `UPDATE "transaction_infos" SET "status"=done WHERE tid = 事务id AND sid = main AND sname = main`
	yoursql.Set_expect(0, 0, sql_key5, []map[string]interface{}{
		{},
	})
	defer yoursql.Clean_expect(sql_key5)

	// 2. do rpc request
	create_token_args := sdk_rpc.Create_token_args{
		Tid:           "事务id",
		Media_type:    1,
		Media_name:    "今日头条",
		Usage:         "测试用",
		App_id:        1234567,
		Secret:        "this is a secret",
		Refresh_token: "xxxrftoken",
	}
	create_token_reply := sdk_rpc.Create_token_reply{}
	err = sdk_rpc.Create_token(rpc_client, &create_token_args, &create_token_reply, 0)
	r.Nil(err)
	r.Equal(create_token_reply.Err.Msg, v1.Err_create_token_ok.Msg)
	r.Equal(sdk_errcode.From_media_api_info_service, create_token_reply.Err.From())
	r.Equal(sdk_errcode.Code_s2s_ok, create_token_reply.Err.Code())
}

/*
func Test_api_info_service_ping(t *testing.T) {
	r := require.New(t)

	// 1. for mock server
	client_conn, s_client_conn := net.Pipe()
	rpc_server := pkg_rpc.New_server()
	service := api_info_service.New()
	err := rpc_server.Register_by_name(service, "v1.api.api_info_service")
	r.Nil(err)
	go rpc_server.Test_json(s_client_conn)

	// 2. for mock client
	rpc_client := jsonrpc.NewClient(client_conn)
	defer rpc_client.Close()
	args := sdk.Ping_args{Msg: "ping"}
	reply := sdk.Ping_reply{}
	err = sdk.Ping(rpc_client, &args, &reply, 0)
	r.Nil(err)
	r.Equal("pong", reply.Msg)
}
*/
