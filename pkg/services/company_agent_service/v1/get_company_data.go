package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	pkg_dbs "yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/test_tool"
	sdk_arpc "yunfan/sdk/arpc/company_agent_service"
	sdk_dbs "yunfan/sdk/dbs/company_agent_service"
	sdk_errcode "yunfan/sdk/errcode"
)

var (
	Err_get_company_data_ok                 = pkg_errcode.New("get_company_data: ok", me, sdk_errcode.Code_s2s_ok)
	Err_get_company_data_need_ids           = pkg_errcode.New("get_company_data: need ids", me, sdk_errcode.Code_s2s_need_arg)
	Err_get_company_data_disconnect_from_db = pkg_errcode.New("get_company_data: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *Company_agent_service) Get_company_data(rwc io.ReadWriteCloser, args *sdk_arpc.Get_company_data_args, reply *sdk_arpc.Get_company_data_reply) error {
	if len(args.Ids) == 0 {
		reply.Err = Err_get_company_data_need_ids
		return nil
	}

	db, ok := pkg_dbs.Begin_sqlite()
	if !ok {
		reply.Err = Err_get_company_data_disconnect_from_db
		return nil
	}
	defer pkg_dbs.End_sqlite()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	company_data_records := []sdk_dbs.Company_data_info{}
	result := db.WithContext(ctx).Where(args.Ids).Find(&company_data_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("get_company_data: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		return nil
	}

	reply.Company_data_info_list = company_data_records
	reply.Err = Err_get_company_data_ok

	return nil
}

// @Summary 查询公司数据
// @Description 查询公司数据
// @Tags company-agent-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Get_company_data_args true "args"
// @Success 200 {object} sdk_arpc.Get_company_data_reply "reply"
// @Router /v1.company-agent-service.Get_company_data [post]
func (this *Company_agent_service) Get_company_data_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Get_company_data_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Get_company_data_reply{}
	if err := this.Get_company_data(&test_tool.Test_conn{}, args, reply); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := enc.Encode(reply); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buffer.Bytes())
}
