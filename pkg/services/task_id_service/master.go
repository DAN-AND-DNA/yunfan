package main

import (
	"context"
	"net"
	"strings"
	"time"
	pkg_errcode "yunfan/pkg/errcode"
	"yunfan/pkg/services/task_id_service/error_code"
	sdk_errcode "yunfan/sdk/errcode"

	clientv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type master struct {
	cli          *clientv3.Client
	ip           string
	key          string
	ttl          int64
	isMasterNode bool
	revision     int64
	id           clientv3.LeaseID
	isClose      bool
}

var Master_node = &master{}

func Init_master_node(etcdAddr []string, port string, ttl int64) error {

	return init_master_node(etcdAddr, Internal_ip()+":"+port, ttl)
}

func Internal_ip() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

func init_master_node(etcdAddr []string, ip string, ttl int64) error {

	Master_node.key = "/api/task_id"
	Master_node.ttl = ttl
	Master_node.ip = ip

	var err error
	Master_node.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdAddr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {

		log_err := pkg_errcode.New("init_master_node: connect etcd "+err.Error()+" ip:"+ip+" etcdAddr:"+strings.Join(etcdAddr, ","), error_code.Me, sdk_errcode.Code_db_node_internal_error)
		error_code.Print(log_err)
		return err
	}
	go Master_node.cron_ttl()
	return nil
}

func (c *master) cron_ttl() {
	if c == nil {
		panic("InitMasterNode is nil")
	}
	if err := c.apply_master_node(); err != nil {
		panic(err)
	}
	c.watch()
	ticker := time.NewTicker(time.Duration(c.ttl) * time.Second)
	go func() {
		for range ticker.C {
			c.apply_master_node()
		}
	}()
}

func (c *master) watch() {
	go func() {
		watcher := clientv3.NewWatcher(c.cli)
		watchChan := watcher.Watch(context.Background(), c.key, clientv3.WithRev(c.revision+1))
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.DELETE:
					if !c.isClose {
						go c.apply_master_node()
					}
				}
			}

		}
	}()
}

func (c *master) apply_master_node() error {
	if c == nil {
		panic("InitMasterNode is nil")
	}
	lease := clientv3.NewLease(c.cli)
	if !c.isMasterNode {
		txn := clientv3.NewKV(c.cli).Txn(context.TODO())
		grantRes, err := lease.Grant(context.TODO(), c.ttl+1)
		if err != nil {

			log_err := pkg_errcode.New("apply_master_node: New "+err.Error()+" ip:"+c.ip, error_code.Me, sdk_errcode.Code_db_node_internal_error)
			error_code.Print(log_err)
			c.isMasterNode = false
			return err
		}
		c.id = grantRes.ID
		txn.If(clientv3.Compare(clientv3.CreateRevision(c.key), "=", 0)).
			Then(clientv3.OpPut(c.key, c.ip, clientv3.WithLease(grantRes.ID))).
			Else()
		txnResp, err := txn.Commit()
		if err != nil {
			log_err := pkg_errcode.New("apply_master_node: New "+err.Error()+" ip:"+c.ip, error_code.Me, sdk_errcode.Code_db_node_internal_error)
			error_code.Print(log_err)
			c.isMasterNode = false
			return err
		}
		if txnResp.Succeeded {
			c.isMasterNode = true
		} else {
			c.isMasterNode = false
		}
		if c.revision != txnResp.Header.Revision {
			c.revision = txnResp.Header.Revision
		}
	}
	_, err := lease.KeepAliveOnce(context.TODO(), c.id)
	if err != nil {
		log_err := pkg_errcode.New("apply_master_node: Keep alive "+err.Error()+" ip:"+c.ip, error_code.Me, sdk_errcode.Code_db_node_internal_error)
		error_code.Print(log_err)
		c.isMasterNode = false
		return err
	}
	return nil
}

func (c *master) Close_apply_master_node() {
	if c != nil {
		c.isClose = true
		if _, err := c.cli.Delete(context.Background(), c.key); err != nil {
			log_err := pkg_errcode.New("apply_master_node: Delete "+err.Error()+" ip:"+c.ip, error_code.Me, sdk_errcode.Code_db_node_internal_error)
			error_code.Print(log_err)
		}
	}
}
