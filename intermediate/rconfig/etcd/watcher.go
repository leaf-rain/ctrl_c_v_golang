package etcd

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	etcdConfig *etcdConfig
	ch         clientv3.WatchChan
	closeChan  chan struct{}
}

func newWatcher(s *etcdConfig) *watcher {
	w := &watcher{
		etcdConfig: s,
		ch:         nil,
		closeChan:  make(chan struct{}),
	}
	w.ch = s.client.Watch(s.ctx, s.options.root, clientv3.WithPrefix())
	return w
}

func (w *watcher) Watcher() {
	for {
		select {
		case <-w.etcdConfig.ctx.Done():
			return
		case <-w.closeChan:
			return
		case kv, ok := <-w.ch:
			if !ok {
				return
			}
			var storgeData []*mvccpb.KeyValue
			var delData []*mvccpb.KeyValue
			for _, v := range kv.Events {
				if v.Type == mvccpb.DELETE {
					delData = append(delData, v.Kv)
				} else {
					storgeData = append(storgeData, v.Kv)
				}
			}
			if len(storgeData) > 0 {
				w.etcdConfig.storge(storgeData)
			}
			if len(delData) > 0 {
				w.etcdConfig.del(delData)
			}
		}
	}
}

func (w *watcher) Close() {
	close(w.closeChan)
	close(w.etcdConfig.reloadChan)
}
