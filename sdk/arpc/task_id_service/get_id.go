package task_id_service

import (
	"context"
	"time"
	"yunfan/pkg/errcode"

	pkg_errcode "yunfan/pkg/errcode"
	sdk_errcode "yunfan/sdk/errcode"

	"snk.git.node1/yunfan/arpc"
)

type Get_type_enum uint8

const (
	Get_empty Get_type_enum = iota
	Get_company
	Get_user
)

type Get_id_args struct {
	Biz_tag string `json:"biz_tag" example:"分类tag"`
}

type Get_id_reply struct {
	Id  int64 `json:"id" example:"1"` // 产生的id
	Err *errcode.Errcode
}

func Init_arpc_client(etcd_addr []string) (*arpc.Client, error) {

	//is mock
	if len(etcd_addr) == 0 {
		return nil, nil
	}

	err := default_client.init_etcd_rpc(etcd_addr, 30)
	if err != nil {
		return nil, err
	}
	c, err := default_client.get_rpc_client()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Refresh_client() (*arpc.Client, error) {
	default_client.m.Lock()
	default_client.change = true
	default_client.m.Unlock()

	c, err := default_client.get_rpc_client()
	if err != nil {
		return nil, err
	}

	return c, nil
}
func Get_id(rpc_client *arpc.Client, args *Get_id_args, reply *Get_id_reply, timeout time.Duration) error {
	//is mock
	if rpc_client == nil {
		get_id_mock(reply)
		return nil
	}
	svc_method := "v1.task-id-service.Get_id"

	if int64(timeout) == 0 {
		return rpc_client.Call(svc_method, args, reply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case call := <-rpc_client.Go(svc_method, args, reply, nil).Done:

		return call.Error
	case <-ctx.Done():
		return ctx.Err()
	}

}

func get_id_mock(reply *Get_id_reply) {
	reply.Id = 1
	reply.Err = pkg_errcode.New("get_id: ok", sdk_errcode.From_task_id_service, sdk_errcode.Code_s2s_ok)
}
