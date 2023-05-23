package rredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

/* 该延时队列为简单实现，暂不支持重试和死信 */

type DelayedQueue interface {
	// 添加队列消息
	Put(ctx context.Context, z []*redis.Z) error
	// 处理单条队列消息
	DealWithOnce(ctx context.Context, min, max float64) (*redis.Z, error)
	// 处理多条队列消息
	DealWithMultiple(ctx context.Context, min, max float64, count int64) ([]redis.Z, error)
	// 队列消息加速
	Expedite(ctx context.Context, member string, score float64) error
}

type QueueDealWith interface {
	DealWithOnce(z redis.Z) error
	DealWithMultiple(z []redis.Z) error
}

type defaultQueue struct {
	keyName  string
	rd       *Client
	lock     *RedisLock
	dealWith QueueDealWith
}

func NewDelayedQueue(rd *Client, fun QueueDealWith, keyName, lockKey, lockVal string, lockTtl time.Duration) DelayedQueue {
	return &defaultQueue{
		keyName:  keyName,
		rd:       rd,
		lock:     NewRedisLock(rd, lockKey, lockVal, lockTtl),
		dealWith: fun,
	}
}

func (d defaultQueue) Put(ctx context.Context, z []*redis.Z) error {
	err := d.rd.ZAdd(ctx, d.keyName, z...).Err()
	if err != nil {
		return err
	}
	return err
}

func (d defaultQueue) Expedite(ctx context.Context, member string, score float64) error {
	return d.rd.ZIncrBy(ctx, d.keyName, score, member).Err()
}

func (d defaultQueue) DealWithOnce(ctx context.Context, min, max float64) (*redis.Z, error) {
	lock, err := d.lock.TryLock(ctx)
	if err != nil {
		return nil, err
	}
	if !lock {
		return nil, nil
	}
	defer d.lock.UnLock(ctx)
	var z []redis.Z
	z, err = d.rd.ZRangeByScoreWithScores(ctx, d.keyName, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", min),
		Max:    fmt.Sprintf("%f", max),
		Offset: 0,
		Count:  1,
	}).Result()
	if len(z) == 0 {
		return nil, nil
	}
	err = d.dealWith.DealWithOnce(z[0])
	if err == nil {
		return &z[0], d.rd.ZRem(ctx, d.keyName, z[0].Member).Err()
	} else {
		return nil, err
	}
}

func (d defaultQueue) DealWithMultiple(ctx context.Context, min, max float64, count int64) ([]redis.Z, error) {
	lock, err := d.lock.TryLock(ctx)
	if err != nil {
		return nil, err
	}
	if !lock {
		return nil, nil
	}
	defer d.lock.UnLock(ctx)
	var z []redis.Z
	z, err = d.rd.ZRangeByScoreWithScores(ctx, d.keyName, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", min),
		Max:    fmt.Sprintf("%f", max),
		Offset: 0,
		Count:  count,
	}).Result()
	if len(z) == 0 {
		return nil, nil
	}
	err = d.dealWith.DealWithMultiple(z)
	if err == nil {
		var member = make([]interface{}, len(z))
		for i := range z {
			member[i] = z[i].Member
		}
		return z, d.rd.ZRem(ctx, d.keyName, member...).Err()
	} else {
		return nil, err
	}
}
