package rredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisSubscribe struct {
	conn *redis.Client
}

func NewRedisSubscribe(conn *redis.Client) *RedisSubscribe {
	return &RedisSubscribe{conn: conn}
}

func (r *RedisSubscribe) PubMessage(ctx context.Context, channelName, msg string) {
	r.conn.Publish(ctx, channelName, msg)
}

func (r *RedisSubscribe) SubMessage(ctx context.Context, channelName string, msgChan chan string) {
	pubsub := r.conn.Subscribe(ctx, channelName)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		msgChan <- msg.Payload
	}
}
