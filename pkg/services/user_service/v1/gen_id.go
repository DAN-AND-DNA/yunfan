package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
	"yunfan/pkg/dbs"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/test_tool"
	sdk_arpc "yunfan/sdk/arpc/user_service"
	sdk_dbs "yunfan/sdk/dbs/user_service"
	sdk_errcode "yunfan/sdk/errcode"

	json "github.com/json-iterator/go"
)

var (
	Err_gen_id_ok                 = pkg_errcode.New("gen_id: ok", me, sdk_errcode.Code_s2s_ok)
	Err_gen_id_need_tid           = pkg_errcode.New("gen_id: need tid", me, sdk_errcode.Code_s2s_need_arg)
	Err_gen_id_need_type          = pkg_errcode.New("gen_id: need type", me, sdk_errcode.Code_s2s_need_arg)
	Err_gen_id_bad_type           = pkg_errcode.New("gen_id: bad type", me, sdk_errcode.Code_s2s_bad_arg)
	Err_gen_id_no_such_trans      = pkg_errcode.New("gen_id: no such trans", me, sdk_errcode.Code_s2s_bad_arg)
	Err_gen_id_disconnect_from_db = pkg_errcode.New("gen_id: disconnect from db", me, sdk_errcode.Code_db_internal_error)
)

func (this *User_service) Gen_id(rwc io.ReadWriteCloser, args *sdk_arpc.Gen_id_args, reply *sdk_arpc.Gen_id_reply) error {
	sid := strconv.Itoa(sdk_errcode.From_gen_id)
	sname := "gen_id"

	// 1. 检查参数
	if args.Tid == "" {
		reply.Err = Err_gen_id_need_tid
		return nil
	}

	switch args.Type {
	case sdk_arpc.Gen_empty:
		reply.Err = Err_gen_id_need_type
		return nil
	case sdk_arpc.Gen_company:
	case sdk_arpc.Gen_user:
	default:
		reply.Err = Err_gen_id_bad_type
		return nil
	}

	// 2. 验证事务和分支
	db, ok := dbs.Postgres()
	if !ok {
		reply.Err = Err_gen_id_disconnect_from_db
		Print(reply.Err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	trans_records := []sdk_dbs.Transaction_info{}
	result := db.WithContext(ctx).Where("tid = ?", args.Tid).Find(&trans_records)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("gen_id: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	has_registerd := false
	has_sub_done := false
	str_id := ""

	if len(trans_records) > 0 {
		has_registerd = true
	}

	if !has_registerd {
		reply.Err = Err_gen_id_no_such_trans
		return nil
	}

	for _, record := range trans_records {
		if record.Sid == sid {
			has_sub_done = true
			str_id = record.Status
			break
		}
	}

	if has_sub_done {
		reply.Id, _ = strconv.ParseUint(str_id, 10, 64)
		reply.Err = Err_gen_id_ok
		return nil
	}

	// 3. 获得id
	now_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), this.loc)
	new_id := (uint64)(0)

	if args.Type == sdk_arpc.Gen_company {
		new_id_record := sdk_dbs.Company_info_id{
			Create_timestamp: (uint64)(now_time.Unix()),
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		result := db.WithContext(ctx).Create(&new_id_record)
		cancel()

		if result.Error != nil {
			reply.Err = pkg_errcode.New("gen_id: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		} else {
			new_id = new_id_record.Id
		}

	} else if args.Type == sdk_arpc.Gen_user {
		new_id_record := sdk_dbs.User_info_id{
			Create_timestamp: (uint64)(now_time.Unix()),
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		result := db.WithContext(ctx).Create(&new_id_record)
		cancel()

		if result.Error != nil {
			reply.Err = pkg_errcode.New("gen_id: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
			Print(reply.Err)
			return nil
		} else {
			new_id = new_id_record.Id
		}
	}

	// 4. 一次完成: 创建公司，确认分支
	new_t_record := &sdk_dbs.Transaction_info{
		Tid:              args.Tid,
		Sid:              sid,
		Sname:            sname,
		Status:           strconv.FormatUint(new_id, 10),
		Create_timestamp: (uint64)(now_time.Unix()),
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	result = db.WithContext(ctx).Create(new_t_record)
	cancel()
	if result.Error != nil {
		reply.Err = pkg_errcode.New("gen_id: "+result.Error.Error(), me, sdk_errcode.Code_db_internal_error)
		Print(reply.Err)
		return nil
	}

	reply.Err = Err_gen_id_ok
	return nil
}

// @Summary 创建id
// @Description 创建id
// @Tags user-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Gen_id_args true "args"
// @Success 200 {object} sdk_arpc.Gen_id_reply "reply"
// @Router /v1.user-service.Gen_id [post]
func (this *User_service) Gen_id_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Gen_id_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Gen_id_reply{}
	if err := this.Gen_id(&test_tool.Test_conn{}, args, reply); err != nil {
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
