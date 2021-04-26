package injection

import (
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/dao"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth/login_check"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

func Injection() {
	task.SetBackupDao(dao.BackupDAO)
	task.SetTaskDao(dao.Task)
	task.SetSubTaskDao(dao.SubTask)
	task.SetDbTools(db_info.DBTool)
	task.SetChecker(checker.Checker)
	checker.SetRuleStatusDao(dao.Rule)
	db_info.SetClusterDao(dao.Cluster)
	auth.SetLoginService(login_check.LoginService)
	admin.SetAdminDao(dao.Admin)

	switch config.Conf.Role.From {
	case "conf":
		task.SetAuthTools(auth.ConfAuthService)
	case "net":
		task.SetAuthTools(auth.NetAuthService)
	default:
		MockInjection()
	}
}

func MockInjection() {
	task.SetAuthTools(auth.MockAuthService)
}
