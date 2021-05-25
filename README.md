### db injection
> sql check, sql exec, backup and recover dml handle, etc.


#### 技术栈
* gin--web框架
* gorm--orm
* ldap--认证

## 系统构建

获取依赖包
```
go mod tidy
```

## 如何启动
config 包下确定配置文件名是否正确 (测试使用 config-example.yml) :
```go
const (
	configPathEnv = "config_path"
	configPath    = "./config/config.yml"
)
```

修改配置文件 config.yml 中 db 对象：
```go
db:
  address: "127.0.0.1"
  port: 3306
  user: "root"
  password: "123456"
  db_name: "dbinjection"
  max_idle_conn: 2
  max_open_conn: 30

```

启动 main.go
```go
func main() {
	flag.Parse()
	log.Println("version:", Version)
	config.InitConfig("")
	logger.InitLog(config.Conf.Server.LogDir, "dbinjection.log", config.Conf.Server.LogLevel)
	dao.InitDB()
	injection.Injection()
	checker.InitRuleStatus()

	router.Run()
}

```

## 测试规则

```go
package test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"
	task "gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

func TestListRule(t *testing.T) {
	SQLContent := "select * from table_a"

	cluster := db_info.DbInjectionCluster{
		ID:          0,
		Name:        "dbinjection",
		Description: "1",
		Addr:        "127.0.0.1",
		User:        "root",
		Pwd:         "123456",
		Ct:          0,
		Ut:          0,
		Operator:    "1",
	}

	var dbName string = "dbinjection"
	defaultDBName := dbName

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, dbName))
	if err != nil {
		fmt.Printf("open db conn err: %s", err.Error())
	}

	defaultDB, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, defaultDBName))
	if err != nil {
		fmt.Printf("open db conn err: %s", err.Error())
	}

	dbInfo := task.DBInfo{DB: db, DefaultDB: defaultDB, DBName: dbName}
	pass, suggestion, affectRow := checker.Checker.SqlCheck(SQLContent, "", "", &dbInfo)
	if affectRow > 0 {
		fmt.Printf("PASS = %t %s %d", pass, suggestion, affectRow)
	}
	if !pass {
		fmt.Printf("PASS = %t %s %d", pass, suggestion, affectRow)
	}

	dbInfo.CloseConn()

}