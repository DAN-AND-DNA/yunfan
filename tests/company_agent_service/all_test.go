package company_agent_service

import (
	"context"
	"net"
	"testing"
	"time"
	"yunfan/pkg/rpc"
	pkg_rpc "yunfan/pkg/rpc"
	v1 "yunfan/pkg/services/company_agent_service/v1"
	sdk_rpc "yunfan/sdk/arpc/company_agent_service"

	pkg_dbs "yunfan/pkg/dbs"

	sdk_dbs "yunfan/sdk/dbs/company_agent_service"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
	"snk.git.node1/yunfan/arpc"
)

var (
	c_conn, s_conn              = net.Pipe()
	pool_media_api_info_service rpc.Std_rpc_client_pool
)

func init() {
	// 1. mock db
	pkg_dbs.Register_sqlite("transaction_info", "", &sdk_dbs.Transaction_info{})
	pkg_dbs.Register_sqlite("company_data_info", "", &sdk_dbs.Company_data_info{})

	if err := pkg_dbs.Open_sqlite("dsn", true); err != nil {
		panic(err)
	}

	// 2. init target rpc service
	service := v1.New()
	rpc_server := pkg_rpc.New_server()

	if err := rpc_server.Register_by_name(service, "v1.company-agent-service"); err != nil {
		panic(err)
	}

	go rpc_server.Test_json(s_conn)

	var err error
	pool_media_api_info_service, err = pkg_rpc.New_json_rpc_client_pool(1, 1,
		// 1. conn_factory
		func() (net.Conn, error) {
			return c_conn, nil
		},

		// 2. keep_alive_callback
		func(rpc_client *arpc.Client) bool {
			return true
		}, 15)

	if err != nil {
		panic(err)
	}
}

func Test_company_agent_service_Get_company_data(t *testing.T) {
	r := require.New(t)

	// 1. insert mock data
	db, ok := pkg_dbs.Begin_sqlite()
	r.Equal(ok, true)
	defer pkg_dbs.End_sqlite()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fake_data_record := sdk_dbs.Company_data_info{
		Id:   123456789,
		Name: "测试公司",
		Cost: 10300.7,
	}

	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&fake_data_record)
	cancel()
	r.Nil(result.Error)

	// 2. do rpc request here
	get_company_data_args := sdk_rpc.Get_company_data_args{Ids: []uint64{123456789}}
	get_company_data_reply := sdk_rpc.Get_company_data_reply{}

	rpc_client, err := pool_media_api_info_service.Get()
	r.Nil(err)
	defer rpc_client.Close()

	err = sdk_rpc.Get_company_data(rpc_client.Client, &get_company_data_args, &get_company_data_reply, 5*time.Second)
	r.Nil(err)
	r.NotNil(get_company_data_reply.Err)
	r.Equal(get_company_data_reply.Err.Code(), 0)
	r.Equal(len(get_company_data_reply.Company_data_info_list), 1)
	r.Equal(get_company_data_reply.Company_data_info_list[0].Id, uint64(123456789))
	r.Equal(get_company_data_reply.Company_data_info_list[0].Name, "测试公司")
	r.Equal(get_company_data_reply.Company_data_info_list[0].Id, uint64(123456789))
}

func Test_company_agent_service_reopen_for_clean(t *testing.T) {
	r := require.New(t)

	err := pkg_dbs.Reopen_sqlite()
	r.Nil(err)

	// 1. insert mock data
	db, ok := pkg_dbs.Begin_sqlite()
	r.Equal(ok, true)
	defer pkg_dbs.End_sqlite()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fake_data_record := sdk_dbs.Company_data_info{
		Id:   123456789,
		Name: "测试公司123",
		Cost: 10300.7,
	}

	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&fake_data_record)
	cancel()
	r.Nil(result.Error)

	// 2. do rpc request here
	get_company_data_args := sdk_rpc.Get_company_data_args{Ids: []uint64{123456789}}
	get_company_data_reply := sdk_rpc.Get_company_data_reply{}

	rpc_client, err := pool_media_api_info_service.Get()
	r.Nil(err)
	defer rpc_client.Close()

	err = sdk_rpc.Get_company_data(rpc_client.Client, &get_company_data_args, &get_company_data_reply, 5*time.Second)
	r.Nil(err)
	r.NotNil(get_company_data_reply.Err)
	r.Equal(get_company_data_reply.Err.Code(), 0)
	r.Equal(len(get_company_data_reply.Company_data_info_list), 1)
	r.Equal(get_company_data_reply.Company_data_info_list[0].Id, uint64(123456789))
	r.Equal(get_company_data_reply.Company_data_info_list[0].Name, "测试公司123")
	r.Equal(get_company_data_reply.Company_data_info_list[0].Id, uint64(123456789))
}
