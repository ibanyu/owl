package task

import "database/sql"

type DBInfo struct {
	DB        *sql.DB
	DefaultDB *sql.DB
	DBName    string
}

type dbTools interface {
	GetDBConn(dbName, cluster string) (*DBInfo, error)
}

var dbTool dbTools

func setDbTools(impl dbTools) {
	dbTool = impl
}

type DbInjectionAccount struct {
	ID       uint64 `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Passwd   string `json:"passwd" gorm:"column:passwd"`
	Dbtype   int64  `json:"dbtype" gorm:"column:dbtype"`
}
