<div align="left">

**简体中文 | [English](../README.md)**

</div>

#### Owl 简介

Owl 是伴鱼基于公司内部场景，开发的一个数据库SQL审核平台，致力于规范研发同学的建表、索引、数据安全等操作。
提供的功能主要包括：
* 登陆认证：基于Ldap的登陆认证
* 流程审批：支持两级审核，管理员(DBA)和研发
* 规则列表：支持动态开启与关闭审核规则
* SQL自动检测：通过TiDB Parser做SQL语法检测
* SQL规则引擎：通过一套规则，规范SQL上线
* SQL定时执行：风险操作，调度到低峰期执行
* 数据备份和回滚：对于DML操作，备份操作前的数据，并支持回滚

![architecture](./image/architecture.png)

#### 功能简介

具体请参考[功能介绍](./introduction.md)

## 部署及使用

### 部署环境

Owl是前后端分离的架构，后端基于go的gin web框架编写，依赖的基础环境包括：
* go 1.3+
* tidb、mysql（主要用于后端元数据存储）

前端基于react的ant design pro框架编写，依赖的基础环境包括：
* node
* yarn (npm)

### 后端单独部署

1、安装go环境
```
[root@dongfengtest-host-0 local]# go version
go version go1.16.7 linux/amd64
```
2、下载后端代码到本地目录
```
[root@dongfengtest-host-0 sql_audit]# git clone https://github.com/ibanyu/owl.git                                                        ^C
[root@dongfengtest-host-0 sql_audit]#
[root@dongfengtest-host-0 sql_audit]# ls
owl
```
3、编辑配置文件，将config/目录下config-example.yml重命名成config.yml，并配置好数据库和Ldap配置
```
db:
  address: "xx.xx.xx.xx"
  port: xx
  user: "xx"
  password: "xx"
  db_name: "owl"
  max_idle_conn: 2
  max_open_conn: 30

login:
  ldap:
    host: "ldap.test.com"
    port: 389
    base_dn: "dc=test,dc=com"
    use_ssl: false
    bind_dn: "cn=hi,dc=test,dc=com"
    bind_pwd: "password"
  login_path: "https://xx.com/login"
  token_secret: ""
  token_effective_hour: 1
```
4、初始化数据库  
* 创建数据库：``` CREATE DATABASE `owl`; use owl ```  
* 初始化表：使用[build_table.sql](../dao/build_table.sql)的sql初始化表  
* 添加首位管理员： ``` insert into owl_admin (username,description) values ('your ldap name','first admin'); ```  

5、编译运行
```
[root@dongfengtest-host-0 owl]# go build -o bin/owl -a ./cmd/owl/
[root@dongfengtest-host-0 owl]#
[root@dongfengtest-host-0 owl]#
[root@dongfengtest-host-0 owl]# ./bin/owl

[2021-08-20 12:50:11] [info] replacing callback `gorm:update_time_stamp` from /data/sql_audit/owl/dao/init.go:36

(/data/sql_audit/owl/dao/rule.go:15)
[2021-08-20 12:50:11]  [1.29ms]  SELECT * FROM `owl_rule_status`
[0 rows affected or returned ]
{"level":"info","ts":"2021-08-20 12:50:11.184","caller":"router/router.go:85","msg":"current dir is: /data/sql_audit/owl/bin"}
{"level":"info","ts":"2021-08-20 12:50:11.184","caller":"router/router.go:111","msg":"start listening port: 8081"}
```

### 前端单独部署

1、前端地址  

[owl_web](https://github.com/ibanyu/owl_web)

2、安装node
```
[root@dongfengtest-host-0 local]# node -v
v16.7.0
```
3、下载前端代码到本地目录
```
git clone https://github.com/ibanyu/owl_web.git
```
4、进入owl_web目录，安装依赖并编译运行
```
bogon:owl_web liujiang$ npm install
bogon:owl_web liujiang$ 
bogon:owl_web liujiang$ vim config/proxy.js  配置后端访问地址
bogon:owl_web liujiang$ 
bogon:owl_web liujiang$ npm start
  App running at:
  - Local:   http://localhost:8000 (copied to clipboard)
  - Network: http://xx.xx.xx.xx:8000
```

### 前后端混合部署或者容器化部署

```
# 仅构建后端
make build

# 交叉编译后端
make build-linux

# 构建UI并置于static目录
make build-front

# 启动; 如需同时启动UI，则需先执行make build-front
make run

# 编译docker镜像; 如需同时启动UI，则需先执行make build-front
make build-docker

# docker 启动
make run-docker
```


#### 如何参与贡献

* 提交bug：通过issue的方式提交bug。
* 贡献代码：fork仓库，进行代码变更，提交pr。
* 参与设计讨论：通过issue或者pr（/doc）提交设计文档，并讨论。

## License

[Apache 2.0 License](./LICENSE)