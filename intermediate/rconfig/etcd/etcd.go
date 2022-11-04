package etcd

import (
	"context"
	"errors"
	"github.com/leaf-rain/ctrl_c_v_golang/intermediate/recuperate"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
)

var (
	ErrClientNil = errors.New("etcd client is nil")
)

type etcdConfig struct {
	ctx        context.Context
	client     *clientv3.Client
	options    *options
	cache      *sync.Map
	once       *sync.Once
	cancel     func()
	watcher    *watcher
	reloadChan chan string
	isWatch    bool
}

func NewEtcdConfig(ctx context.Context, cli *clientv3.Client, opts ...Option) (*etcdConfig, error) {
	if cli == nil {
		return nil, ErrClientNil
	}
	var cancel func()
	if ctx == nil {
		ctx, cancel = context.WithCancel(context.Background())
	}
	conf := &etcdConfig{
		ctx:     ctx,
		client:  cli,
		options: NewOptions(opts...),
		cache:   new(sync.Map),
		once:    new(sync.Once),
		cancel:  cancel,
	}
	conf.watcher = newWatcher(conf)
	if err := conf.Load(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *etcdConfig) Load() error {
	for _, v := range c.options.paths {
		if err := c.loadPath(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *etcdConfig) Get() *sync.Map {
	return c.cache
}

func (c *etcdConfig) loadPath(path string) error {
	rsp, err := c.client.Get(c.ctx, c.options.root+path)
	if err != nil {
		return err
	}
	for _, item := range rsp.Kvs {
		k := string(item.Key)
		k = strings.ReplaceAll(k, c.options.root, "")
		if k == path {
			c.cache.Store(k, item.Value)
			break
		}
	}
	return nil
}

func (c *etcdConfig) getPath(str string) string {
	return strings.ReplaceAll(str, c.options.root, "")
}

func (c *etcdConfig) storge(kvs []*mvccpb.KeyValue) {
	for _, item := range kvs {
		k := string(item.Key)
		k = c.getPath(k)
		for _, v := range c.options.paths {
			if k == v {
				c.cache.Store(k, item.Value)
				if c.isWatch {
					c.reloadChan <- k
				}
			}
		}
	}
}

func (c *etcdConfig) del(kvs []*mvccpb.KeyValue) {
	for _, item := range kvs {
		k := string(item.Key)
		k = c.getPath(k)
		for _, v := range c.options.paths {
			if k == v {
				c.cache.Delete(k)
				if c.isWatch {
					c.reloadChan <- k
				}
			}
		}
	}
}

func (c *etcdConfig) Watch() <-chan string {
	c.once.Do(func() { // 只执行一次
		c.isWatch = true // 只有当调用watch以后才开始监听
		c.reloadChan = make(chan string)
		go func() {
			defer recuperate.Recover()
			c.watcher.Watcher()
		}()
	})
	return c.reloadChan
}

func (c *etcdConfig) Close() {
	c.watcher.Close()
}
