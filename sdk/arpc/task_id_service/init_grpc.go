package task_id_service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"snk.git.node1/yunfan/arpc"
)

type Client struct {
	etcd      *clientv3.Client
	node      string
	key       string
	change    bool
	ttl       int64
	revision  int64
	conn      *arpc.Client
	is_active bool
	m         sync.Mutex
}

var (
	default_client = &Client{}
)

func (cli *Client) init_etcd_rpc(etcdAddr []string, ttl int64) error {
	if cli.is_active {
		return nil
	}

	cli.m.Lock()
	defer cli.m.Unlock()

	c, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	cli.etcd = c
	cli.change = false
	cli.key = "/api/task_id"
	cli.ttl = ttl
	cli.cornTTL()
	cli.is_active = true
	return nil
}

func (c *Client) watch() {
	watcher := clientv3.NewWatcher(c.etcd)
	watchChan := watcher.Watch(context.Background(), c.key, clientv3.WithRev(c.revision+1))
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.DELETE:
				go c.getMasterNode()
			}
		}
	}
}

func (c *Client) cornTTL() {
	if err := c.getMasterNode(); err != nil {
		panic(err)
	}
	go c.watch()
	ticker := time.NewTicker(time.Duration(c.ttl) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = c.getMasterNode()
			}
		}
	}()
}

func (c *Client) getMasterNode() error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := c.etcd.Get(ctx, c.key)
	if err != nil {
		return err
	}
	if len(res.Kvs) == 0 {
		return errors.New("etcd key is invalid")
	}
	for _, v := range res.Kvs {
		val := v
		if string(val.Key) == c.key {
			newNode := string(val.Value)
			if c.node != newNode {
				c.change = true
			}
			c.node = string(val.Value)
		}
	}
	if c.revision != res.Header.Revision {
		c.revision = res.Header.Revision
	}
	return nil
}

func (c *Client) get_rpc_client() (*arpc.Client, error) {
	var err error
	if c.change {
		c.m.Lock()
		defer c.m.Unlock()
		c.conn, err = arpc.Dial("tcp", c.node)
		if err != nil {
			return nil, err
		}
		c.change = false
	}
	return c.conn, nil
}
