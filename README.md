# 项目名称
空调项目

## 主要文件夹含义
1. cmd 二进制程序编译入口
3. conf 配置文件存储目录
4. daemon 调用第三方接口、模型端交互
4. router 前端交互

## 编译方式
```
GOARCH=amd64 GOOS=linux go build -o main ./cmd
```

## 增加字段需要修改

1、model 文件夹里面对象

2、数据库

3、一些需要处理的字段（求平均值等）

## 用其他字段过滤

4、数据过滤 filter(jsonData2["Pd"]) 

## 多少个点位

5、point > 24 

## 多少个组成一条数据

6、batchSize := 30 多少条一个批次

## 获取token

token.go 文件里面的url

## 室内机、室外机数量更改

daemon.go  outDeviceUrlMap innerDeviceUrlMap

## 文档相关

待补充


## 如何贡献
贡献patch流程、质量要求

## 讨论
讨论群：XXX

