package rconfig

import (
	"context"
	"github.com/leaf-rain/ctrl_c_v_golang/intermediate/recuperate"
	"log"
	"sync"
)

type Rconfig interface {
	// Get 获取值(每次使用都要重新序列化)
	Get(path string) Val
	// MapValue 获取值(已经序列化好内容)
	MapValue(path, key string) interface{}
	// Map 获取值(已经序列化好内容)
	Map(path string) map[interface{}]interface{}
	// Watch 监听
	Watch()
	// Close 关闭
	Close()
}

const (
	js = "json"
	ya = "yaml"
)

func NewRConfig(watcher Watcher, way string) (Rconfig, error) {
	var conf = &rconfig{
		parseWay: way,
		watcher:  watcher,
		cache:    new(sync.Map),
		once:     new(sync.Once),
	}
	watcher.Get().Range(func(key, value interface{}) bool {
		_ = conf.mapDataToCache(key.(string))
		return true
	})
	return conf, nil
}

type rconfig struct {
	parseWay string
	ctx      context.Context
	watcher  Watcher
	cache    *sync.Map
	once     *sync.Once
}

func (r *rconfig) mapDataToCache(path string) error {
	p, ok := r.watcher.Get().Load(path)
	if !ok {
		r.cache.Delete(path)
		return nil
	}
	var result = make(map[interface{}]interface{})
	var err error
	switch r.parseWay {
	case js:
		err = Val(p.([]byte)).FormatJson(&result)
		if err != nil {
			return err
		}
	case ya:
		err = Val(p.([]byte)).FormatYaml(&result)
		if err != nil {
			return err
		}
	}
	r.cache.Store(path, result)
	return nil
}

func (r *rconfig) MapValue(path, key string) interface{} {
	p, ok := r.cache.Load(path)
	var data map[interface{}]interface{}
	var result interface{}
	if ok {
		data, ok = p.(map[interface{}]interface{})
		if ok {
			result, ok = data[key]
			if ok {
				return result
			}
		}
	}
	return nil
}

func (r *rconfig) Map(path string) map[interface{}]interface{} {
	p, ok := r.cache.Load(path)
	var data map[interface{}]interface{}
	if ok {
		data, ok = p.(map[interface{}]interface{})
		if ok {
			return data
		}
	}
	return nil
}

func (r *rconfig) Get(path string) Val {
	p, ok := r.watcher.Get().Load(path)
	if ok {
		return p.([]byte)
	}
	return nil
}

func (r *rconfig) Watch() {
	r.once.Do(func() {
		go func() {
			defer recuperate.Recover()
			ch := r.watcher.Watch()
			var err error
			for {
				select {
				case data, ok := <-ch:
					if !ok {
						return
					}
					err = r.mapDataToCache(data)
					if err != nil {
						log.Printf("[mapDataToCache] failed, err:%v", err)
					}
				}
			}
		}()
	})
}

func (r *rconfig) Close() {
	r.watcher.Close()
}
