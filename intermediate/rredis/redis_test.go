package rredis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type testQueue struct {
}

func (t testQueue) DealWithOnce(z redis.Z) error {
	fmt.Println("info --->", z)
	return nil
}

func (t testQueue) DealWithMultiple(z []redis.Z) error {
	fmt.Println("slice --->", z)
	return nil
}

var queueName = "test_queue"

func TestNew(t *testing.T) {
	fmt.Println("1", time.Now().Format("2006-01-02 15:04:05"))
	var ctx = context.Background()
	cli, err := NewRedis(Config{
		PoolSize: 5,
		Addr: []string{
			"127.0.0.1:6379",
		},
		DialTimeout: 10,
	}, ctx)
	if err != nil {
		panic(err)
	}
	var with = testQueue{}
	queue := NewDelayedQueue(cli, with, queueName, "test_lock", "text_lock_val", time.Second*10)
	//err = queue.Put(ctx, []*redis.Z{
	//	{
	//		Score:  float64(time.Now().Unix()),
	//		Member: "test1",
	//	},
	//	{
	//		Score:  float64(time.Now().Unix()),
	//		Member: "test2",
	//	},
	//	{
	//		Score:  float64(time.Now().Unix()),
	//		Member: "test3",
	//	},
	//})
	if err != nil {
		panic(err)
	}
	_, err = queue.DealWithOnce(ctx, 0, 1663926158)
	if err != nil {
		t.Errorf("failed,err:%v", err)
	} else {
		t.Log("success")
	}
}
