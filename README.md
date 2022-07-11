# 日志收集客户端

## 项目说明
通过监听日志文件，追踪最新日志，追踪文件使用了[tailf](https://github.com/hpcloud/tail)包

追踪最新日志后，将日志发送至kafka，使用[sarama](https://github.com/Shopify/sarama)操作kafka

## 项目亮点
集成了`etcd`对需要监听的日志文件做了动态配置，系统使用`watch`监听`etcd`的动态变化，并更新tail任务，可同时监听多个日志文件

### etcd配置示例
`etcd`中保存json格式的配置，目前为`path` 、`topic`字段，示例如下：
```json
[
    {
        "path":"c:/xx.log",
        "topic":"web_log"
    },
    {
        "path":"d:/xx.log",
        "topic":"client_log"
    }
]
```
`etcd`中配置的key格式为`collect_log_%s_conf`,其中占位符为本机的ip，项目启动时会获取本机的ip，向`etcd`更新配置时需要先获取部署机器的ip,也可观察项目启动日志，首行即打印ip信息

## 项目不足
当logagent意外宕机或者配置移除监听后再次监听日志文件时，项目无法回到日志文件上一次位置，只能从文件最新位置开始读取
