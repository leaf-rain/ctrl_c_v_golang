package etcd

import (
	"context"
	"fmt"
	"testing"
)
import clientv3 "go.etcd.io/etcd/client/v3"

func TestNewEtcdConfig(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	c, err := NewEtcdConfig(context.Background(), cli, WithRoot("/key"), WithPaths("k1", "k2.json"))
	if err != nil {
		t.Error(err)
		return
	}
	ch := c.Watch()
	c.Get().Range(func(key, value interface{}) bool {
		fmt.Println("k --->", key.(string))
		fmt.Println("value --->", string(value.([]byte)))
		return true
	})
	for {
		select {
		case path, ok := <-ch:
			if !ok {
				return
			}
			t.Logf("检测到变化，path:%s", path)
			data := c.Get()
			if data == nil {
				panic("nil")
			} else {
				data.Range(func(key, value interface{}) bool {
					fmt.Println("k --->", key.(string))
					fmt.Println("value --->", string(value.([]byte)))
					return true
				})
			}
		}
	}
}
