package rconfig

import (
	"context"
	"github.com/leaf-rain/ctrl_c_v_golang/intermediate/rconfig/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestNewHConfig(t *testing.T) {
	var root = "/ddz"
	var p1 = "agent1"
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	var etcdConfig Watcher
	etcdConfig, err = etcd.NewEtcdConfig(context.Background(), cli, etcd.WithRoot(root), etcd.WithPaths(p1))
	if err != nil {
		t.Error(err)
		return
	}
	var rconf Rconfig
	rconf, err = NewRConfig(etcdConfig, "yaml")
	if err != nil {
		t.Error(err)
		return
	}
	rconf.Watch()
	for {
		t.Logf("p1 ---> %v", rconf.Map(p1))
		time.Sleep(time.Second * 3)
	}
}
