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

具体请参考[功能介绍](./introduction.md)

## 部署及使用

###部署环境
Owl是前后端分离的架构，后端基于go的gin web框架编写，依赖的基础环境包括：
* go 1.3+
* tidb、mysql（主要用于后端元数据存储）

前端基于react的ant design pro框架编写，依赖的基础环境包括：
* node
* yarn (npm)

###后端单独部署
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
```cassandraql
db:
  address: "xx.xx.xx.xx"
  port: xx
  user: "xx"
  password: "xx"
  db_name: "dbinjection"
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
4、db表初始化(创建好后端用到的表)，表创建成功后，并在db_injection_admin添加一条管理员记录
```cassandraql
insert into db_injection_admin (username,description) values ('your ldap name','first admin');
```
5、编译运行
```cassandraql
[root@dongfengtest-host-0 owl]# go build -o bin/owl -a ./cmd/owl/
[root@dongfengtest-host-0 owl]#
[root@dongfengtest-host-0 owl]#
[root@dongfengtest-host-0 owl]# ./bin/owl

[2021-08-20 12:50:11] [info] replacing callback `gorm:update_time_stamp` from /data/sql_audit/owl/dao/init.go:36

(/data/sql_audit/owl/dao/rule.go:15)
[2021-08-20 12:50:11]  [1.29ms]  SELECT * FROM `db_injection_rule_status`
[0 rows affected or returned ]
{"level":"info","ts":"2021-08-20 12:50:11.184","caller":"router/router.go:85","msg":"current dir is: /data/sql_audit/owl/bin"}
{"level":"info","ts":"2021-08-20 12:50:11.184","caller":"router/router.go:111","msg":"start listening port: 8081"}
```
### 前端单独部署
1、安装node
```cassandraql
[root@dongfengtest-host-0 local]# node -v
v16.7.0
```
2、npm安装
```cassandraql
[root@dongfengtest-host-0 local]# npm install -g cnpm --registry=https://registry.npm.taobao.org
[root@dongfengtest-host-0 local]# npm -v
7.20.3
```
3、安装webpack
```cassandraql
[root@dongfengtest-host-0 local]# npm install webpack -g
```
4、安装vue
```cassandraql
[root@dongfengtest-host-0 local]# npm install -g @vue/cli
[root@dongfengtest-host-0 local]# vue -V
@vue/cli 4.5.13
```
5、下载前端代码到本地目录
```cassandraql
git clone https://github.com/ibanyu/owl_web.git
```
6、进入owl目录，安装依赖并编译运行
```cassandraql
bogon:db_injection_web liujiang$ npm install
bogon:db_injection_web liujiang$ 
bogon:db_injection_web liujiang$ vim config/proxy.js  配置后端访问地址
bogon:db_injection_web liujiang$ 
bogon:db_injection_web liujiang$ npm start
  App running at:
  - Local:   http://localhost:8000 (copied to clipboard)
  - Network: http://xx.xx.xx.xx:8000
```

###前后端混合部署或者容器化部署
可以参考Makefile文件
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