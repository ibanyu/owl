package task

import "database/sql"

type sqlChecker interface {
	SqlCheck(sql, charset, collation string, info *sql.DB) (pass bool, suggestion string, affectRow int)
	ListRules() interface{}
}

var checker sqlChecker

func setChecker(impl sqlChecker) {
	checker = impl
}
