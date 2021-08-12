<div align="left">

**简体中文 | [English](../README.md)**

</div>

## dbinjection 简介

dbinjeciotn 是一个数据库管理平台，致力于规范化数据库中的数据、索引、以及对db进行的操作，从而减少风险，规避故障。  
提供的功能包括：流程审批、sql自动审查、sql执行、定时执行、数据变更备份及回滚等。

#### 技术栈
* gin
* gorm
* ldap
* pingcap/parser

* react
* Ant Design of React

#### 功能简介

* 工单审批：对sql工单进行简单的流程审批：研发提交、系统审核、dba审核、执行。
* 系统审核：根据一些预定义的规则审核sql，规则包括：命名规范、不允许控制规范、索引匹配检测、变更影响的数据行数限制等。如有需要，可以关闭部分规则。
* sql执行及定时执行： 根据所管理的db集群去执行sql，并支持定时执行。
* 数据变更备份及回滚：数据变更操作执行后会备份原始数据，可以查看原始数据并执行回滚。

## 开发及部署

#### 依赖
* go环境1.3+
* tidb or mysql

* node 
* yarn (npm)

#### 配置文件

* 配置文件位置可由环境变量（"config_path"）或者main函数中指定，不指定则为默认值："./config/config.yml"。  
* 默认情况下，修改/config/config-example.yml 为 /config/config.yml即可正确读取配置文件  
* 然后还需要修改config.yml中db、当前环境等相关的配置。

#### db初始化

* 建表：client 链接数据，复制/dao/build_table.sql中的内容并执行
* 初始化admin： ``` insert into db_injection_admin (username,description) values ('your ldap name','first admin');```

#### 构建及启动
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

## 快速上手

#### 如何使用

请参考 [快速上手](.....)

#### 如何参与贡献

* 提交bug：通过issue的方式提交bug。
* 贡献代码：fork仓库，进行代码变更，提交pr。
* 参与设计讨论：通过issue或者pr（/doc）提交设计文档，并讨论。

## License

[Apache 2.0 License](./LICENSE)