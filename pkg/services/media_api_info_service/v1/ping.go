package v1

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/test_tool"
	sdk_arpc "yunfan/sdk/arpc/media_api_info_service"
	sdk_errcode "yunfan/sdk/errcode"

	json "github.com/json-iterator/go"
)

var (
	Err_ping_ok       = pkg_errcode.New("ping: ok", me, sdk_errcode.Code_s2s_ok)
	Err_ping_need_msg = pkg_errcode.New("ping: need msg", me, sdk_errcode.Code_s2s_need_arg)
	Err_ping_bad_msg  = pkg_errcode.New("ping: bad msg", me, sdk_errcode.Code_s2s_bad_arg)
)

func (this *Media_api_info_service) Ping(rwc io.ReadWriteCloser, args *sdk_arpc.Ping_args, reply *sdk_arpc.Ping_reply) error {
	if conn, ok := rwc.(net.Conn); ok {
		if args.Msg == "" {
			reply.Err = Err_ping_need_msg
			return nil
		}

		if args.Msg != "ping" {
			reply.Err = Err_ping_bad_msg
			return nil
		}

		conn.SetDeadline(time.Now().Add(30 * time.Second))
		reply.Msg = "pong"
		reply.Err = Err_ping_ok
	}
	return nil
}

// @Summary ping
// @Description ping
// @Tags media-api-info-service
// @Accept json
// @Produce json
// @Param args body sdk_arpc.Ping_args true "args"
// @Success 200 {object} sdk_arpc.Ping_reply "reply"
// @Router /v1.media-api-info-service.Ping [post]
func (this *Media_api_info_service) Ping_swag(w http.ResponseWriter, r *http.Request) {
	var buffer = bytes.Buffer{}
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(&buffer)
	args := &sdk_arpc.Ping_args{}
	if err := dec.Decode(args); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	reply := &sdk_arpc.Ping_reply{}
	if err := this.Ping(&test_tool.Test_conn{}, args, reply); err != nil {
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
