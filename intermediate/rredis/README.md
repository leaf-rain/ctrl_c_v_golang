# redis基础封装*(基于go-redis)
***
ctrl c/v永不过时 ！！！

```text
├─delayedQueue.go          实现了一个基于redis的延时队列(不支持重试和死信)
├─lock.go                  简单实现了一个基于redis的分布式锁方案
├─redis.go                 基于redis配置实现一个可以支持集群或者单点的初始化类
├─script.go                存放了一些在工作用使用到的脚本比如扣减库存
├─subscribe.go             实现了一个简单的redis订阅发布的封装
├─struct_to_map            golang的struct序列化到redis的hash里面(偷懒用，不建议使用，底层使用反射，这玩意懂得都懂)
│    ├─decode.go               解码
│    ├─encode.go               加码
│    ├─map.go                  底层使用反射解析代码(原理其实很简单，基础数据不用管，不是基础数据的可以用proto序列化就用proto，不能的话，直接干json,所以慢而且占内存，但胜在方便)
│    ├─struct_to_map.go        封装出对外使用接口
│    ├─tag.go                  tag解析
```