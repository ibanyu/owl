package injection

import (
	"gitlab.pri.ibanyu.com/middleware/dbinjection/dao"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

func Injection() {
	task.SetBackupDao(dao.BackupDAO)
	task.SetTaskDao(dao.Task)
	task.SetSubTaskDao(dao.SubTask)
	task.SetDbTools(db_info.DBTool)
	checker.SetRuleStatusDao(dao.Rule)
	db_info.SetClusterDao(dao.Cluster)

	MockInjection()
}

func MockInjection() {
	task.SetAuthTools(auth.MockAuthService)
}
