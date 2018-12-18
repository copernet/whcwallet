wallet重写后端框架选型
一、语言和数据库
后端开发语言：golang
数据库：mysql
缓存：redis
二、部署
部署环境：国内阿里云，国外aws
三、golang框架选择如下
web框架：https://github.com/gin-gonic/gin
orm框架：http://gorm.io/docs/index.html
websocket: https://github.com/gorilla/websocket
redis: https://github.com/gomodule/redigo
第一期实现先不用websocket，前端轮询后端，在交易所功能中评估是否上websocket
三、代码框架
1. webapi服务目录如下：
webapi
--cmd
--doc
--src
----main.go
----api
--------tx.go
--------balance.go
----model
--------tx.go
--------balance.go
----db
--------mysql.go
--------redis.go
----rpc
--------rpcclient.go
----util
----config
----log
----route
--------tx.go
--------balance.go
----middle
--------filter.go
----vender



2. engine服务
独立部署go服务，负责解析链上数据到数据库和缓存
3.  static项目由前端维护



