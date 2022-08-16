package v1

import (
	"bytes"
	"io"
	"net/http"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/test_tool"
	sdk_arpc "yunfan/sdk/arpc/user_service"
	sdk_errcode "yunfan/sdk/errcode"

	json "github.com/json-iterator/go"
)

var (
	Err_get_auth_info_ok = pkg_errcode.New("get_auth_info: ok", me, sdk_errcode.Code_s2s_ok)
)

func (this *User_service) Get_auth_info(rwc io.ReadWriteCloser, args *sdk_arpc.Get_auth_info_args, reply *sdk_arpc.Get_auth_info_reply) error {

	reply.Err = Err_get_auth_info_ok
	return nil
}

// @Summary 获得授权信息
// @Description 获得授权信息
// @Tags user-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Get_auth_info_args true "args"
// @Success 200 {object} sdk_arpc.Get_auth_info_reply "reply"
// @Router /v1.user-service.Get_auth_info [post]
func (this *User_service) Get_auth_info_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Get_auth_info_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Get_auth_info_reply{}
	if err := this.Get_auth_info(&test_tool.Test_conn{}, args, reply); err != nil {
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
