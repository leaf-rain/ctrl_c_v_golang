package rconfig

import "sync"

type Watcher interface {
	// Watch 监听
	Watch() <-chan string
	// Get 获取配置
	Get() *sync.Map
	// Close 关闭
	Close()
}
