# Primetheus 一些监控

## 现状

- 已经有了Counter，Guage监控

## 目标

- 添加Histogram、Summery，用于接口耗时的监控

## 方案

- 单个server单独注册metrics


## 增加golib形式

- namespace ,subsystem 注册
- 耗时，请求量，返回码监控