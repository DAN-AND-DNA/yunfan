package media_api_info_service

import (
	"net"
	"testing"
	sdk_rpc "yunfan/sdk/arpc/task_id_service"

	pkg_dbs "yunfan/pkg/dbs"
	pkg_rpc "yunfan/pkg/rpc"
	v1 "yunfan/pkg/services/task_id_service/v1"

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
	service        *v1.Task_id_service
)

func init() {
	// 1. mock db
	if err := pkg_dbs.Open_postgres("dsn", true); err != nil {
		panic(err)
	}

	// 2. init target rpc server
	service = v1.New()
	rpc_server := pkg_rpc.New_server()
	if err := rpc_server.Register_by_name(service, "v1.task-id-service"); err != nil {
		panic(err)
	}
	go rpc_server.Test_json(s_conn)

}

// func Test_media_api_info_service_Ping(t *testing.T) {
// 	r := require.New(t)

// 	ping_args := sdk_rpc.Ping_args{Msg: "ping"}
// 	ping_reply := sdk_rpc.Ping_reply{}
// 	err := sdk_rpc.Ping_raw(rpc_client, &ping_args, &ping_reply, 0)
// 	r.Nil(err)
// 	r.Equal(sdk_errcode.From_media_api_info_service, ping_reply.Err.From())
// 	r.Equal(sdk_errcode.Code_s2s_ok, ping_reply.Err.Code())
// 	r.Equal("pong", ping_reply.Msg)
// }

func Test_task_id_get_id(t *testing.T) {
	r := require.New(t)

	// 1. mock db result
	sql_key1 := `SELECT * FROM "segments"`
	err := yoursql.Set_expect(0, 0, sql_key1, []map[string]interface{}{
		{},
	})
	if err != nil {
		fmt.Println(err)
	}
	defer yoursql.Clean_expect(sql_key1)

	now_time, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	r.Nil(err)

	sql_key2 := fmt.Sprintf(`INSERT INTO "segments" ("biz_tag","max_id","step","create_time","update_time") VALUES (test,1001,1000,%d,%d)`, now_time.Unix(), now_time.Unix())
	yoursql.Set_expect(0, 0, sql_key2, []map[string]interface{}{
		{"id": 37},
	})
	defer yoursql.Clean_expect(sql_key2)

	sql_key3 := `SELECT * FROM "segments" WHERE biz_tag = test ORDER BY "segments"."biz_tag" LIMIT 1`
	yoursql.Set_expect(0, 0, sql_key3, []map[string]interface{}{
		{
			"biz_tag":     "test",
			"max_id":      1001,
			"step":        1000,
			"create_time": now_time.Unix(),
			"update_time": now_time.Unix(),
		},
	})
	//defer yoursql.Clean_expect(sql_key3)

	sql_key4 := fmt.Sprintf(`update segments set max_id=max_id+step,update_time = %d where biz_tag = test`, now_time.Unix())
	yoursql.Set_expect(0, 0, sql_key4, []map[string]interface{}{
		{},
	})
	defer yoursql.Clean_expect(sql_key4)

	//2. do rpc request

	service.Init_service()

	_, err = sdk_rpc.Init_arpc_client(nil)
	r.Nil(err)
	get_id_args := sdk_rpc.Get_id_args{
		Biz_tag: "test",
	}
	get_id_reply := sdk_rpc.Get_id_reply{}

	err = sdk_rpc.Get_id(rpc_client, &get_id_args, &get_id_reply, 0)
	r.Nil(err)
	r.Equal(get_id_reply.Err.Msg, v1.Err_get_id_ok.Msg)
	r.Equal(sdk_errcode.From_task_id_service, get_id_reply.Err.From())
	r.Equal(sdk_errcode.Code_s2s_ok, get_id_reply.Err.Code())
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
