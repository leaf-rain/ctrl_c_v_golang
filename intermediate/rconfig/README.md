# 配置中心(目前是基于etcd)
***
ctrl c/v永不过时 ！！！

>本来是想用[viper](https://github.com/spf13/viper)的，毕竟一看官方文档，
我靠,直接支持远程etcd，这tm不香？结果配半天搞不成，各种环境依赖打架，go mod支持太不友好了。
工具嘛，要尽可能的轻量,我尽可能的设计成插件化，即使不用rconfig，用etcd也可直接使用。

```text
├─bytes.go          用于把字节类型序列化成其它类型
├─rconfig.go        代码入口配置
├─rconfig_test.go   代码单元测试
├─watcher.go        接口，只要实现该接口就就可以使用rconfig，为以后对接apollo或者其它配置中心预留。
│    
├─etcd
│    ├─etcd.go      etcd服务
│    ├─etcd_test.go etcd服务单元测试
│    ├─options.go   选项，用于配置etcd路径
│    └─watcher.go   etcd监测变化方法
```