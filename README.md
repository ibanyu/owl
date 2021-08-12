<div align="left">

**[简体中文](./doc/README-zh-CN.md) | English**

</div>

## What is dbinjection

dbinjeciotn is a db manager platform,committed to standardizing the data, index in the database and operations to the database,  to avoid risks and failures.  
capabilities which dbinjection provide include Process approval、sql Audit、sql execute and execute as crontab 、data backup and recover .

#### Technology stack
* gin
* gorm
* ldap
* pingcap/parser

* react
* Ant Design of React

#### Features

* Process approval：approval or reject sql process order：develop submit、system check、dba check、exec sql.
* sql Audit：check sql by some rules, which is Predefined. There are rules like: standardizing name、'null' value not allowed 、index match check、 data change affect lines limit, and so on. 
* sql execute and execute as crontab： After sql audit and approval, the admin can execute sql, or set a feature time to execute sql.
* data backup and recover：Before the data is changed, it will be backup, and you recover it if something is wrong.

Get Started using dbinjeciotn[quick start](...)

## Develop and deployment

#### Dependency
* go 1.3+
* tidb or mysql

* node 
* yarn (npm)

#### Config file

* config file's location can be set by env-param("config_path") or set in main function. default location is "./config/config.yml".  
* in default value, mv '/config/config-example.yml' to '/config/config.yml' can make config work. 
* at last, still need to change config about database, env and so on.

#### DB init

* create table：use mysql client connect db , copy sql in '/dao/build_table.sql' and execute.
* init first admin： ``` insert into db_injection_admin (username,description) values ('your ldap name','first admin');```

#### Build and run
```
# build back end
make build

# build linux back end
make build-linux

# build front and set to static dir
make build-front

# start; if need UI, exec 'make build-front' first
make run

# build docker image, if  need UI, exec 'make build-front' first
make build-docker

# run as docker container
make run-docker
```

## Getting Started

 [introduction](.....)

## Become a contributor

* Contribute to the codebase.
* Contribute to the docs.
* Report and triage bugs by issue.
* Write technical documentations and blog posts, for users and contributors.

## License

[Apache 2.0 License](./LICENSE)
