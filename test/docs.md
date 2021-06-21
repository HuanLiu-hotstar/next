# Primetheus 一些监控

## 现状

- 已经有了Counter，Guage监控

## 目标

- 添加Histogram、Summery，用于接口耗时的监控
- promauto new 出来的变量自动进行了注册

## 方案

- 单个server单独注册metrics

## 增加golib形式

- namespace ,subsystem 注册
- 耗时，请求量，返回码监控



## 服务监控代码

```golang 
AddRequestLatency(method, path string, statusCode int, latency float64)
webfx.AddRequestLatency("GET", string(path)[0:3]+"/m/:match_id", ctx.Response().GetStatusCode(), time.Since(startTime).Seconds())

```

## Client端主调服务的耗时和返回码监控

主调名字，背调名字，uri，以及耗时，返回码


## 内存定位

```go
引入限流器
控制发送signal速度
14k qps 时，内存1.023G，CPU:44%

signal 消费速度最大17k qps


redis客户端内容：
router->base->RedisReplica 会new 一个默认的redis.Client 分配内存较大
RedisDatabase 会new 一个默认的redis.Client 分配内存较大
unicast.New 会分配两个redis.Client，一个是chatRedis, 一个是brokerRedis



go pprof 命令

go tool pprof 10.140.79.217:1996/debug/pprof/heap

```

## 栈内存监控

```sh
没有任务的时候信息 
# runtime.MemStats
# Alloc = 172396896 =>>>> 172M 
# TotalAlloc = 21723191656
# Sys = 433315928         433M
# Lookups = 0
# Mallocs = 305743347	  305M
# Frees = 305558328 	305M	
# HeapAlloc = 172396896    172M
# HeapSys = 375914496	375M
# HeapIdle = 199221248  199M
# HeapInuse = 176693248  176M
# HeapReleased = 161677312 161M
# HeapObjects = 185019
# Stack = 26738688 / 26738688 ===>>>>  26MB 
# MSpan = 669120 / 2621440
# MCache = 43200 / 49152
# BuckHashSys = 1554733
# GCSys = 19338304
# OtherSys = 7099115
# NextGC = 337898688
# LastGC = 1623405393963487710




有任务之后的
# runtime.MemStats
# Alloc = 205090928
# TotalAlloc = 23205903576
# Sys = 433315928
# Lookups = 0
# Mallocs = 326563161
# Frees = 325934402
# HeapAlloc = 205090928
# HeapSys = 375914496
# HeapIdle = 166256640
# HeapInuse = 209657856
# HeapReleased = 37773312
# HeapObjects = 628759
# Stack = 26738688 / 26738688
# MSpan = 1117376 / 2621440
# MCache = 43200 / 49152
# BuckHashSys = 1558925
# GCSys = 19365072
# OtherSys = 7068155
# NextGC = 333489216
# LastGC = 1623406106463575178

```