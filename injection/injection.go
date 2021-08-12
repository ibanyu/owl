package injection

import (
	"github.com/ibanyu/owl/config"
	"github.com/ibanyu/owl/dao"
	"github.com/ibanyu/owl/service"
	"github.com/ibanyu/owl/service/admin"
	"github.com/ibanyu/owl/service/auth"
	"github.com/ibanyu/owl/service/auth/login_check"
	"github.com/ibanyu/owl/service/checker"
	"github.com/ibanyu/owl/service/db_info"
	"github.com/ibanyu/owl/service/task"
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
	service.SetClock(service.RealClock{})
	admin.SetAdminDao(dao.Admin)

	switch config.Conf.Role.From {
	case "conf":
		task.SetAuthTools(auth.ConfAuthService)
	case "net":
		task.SetAuthTools(auth.NetAuthService)
	case "admin":
		task.SetAuthTools(auth.AdminAuthService)
	case "mock":
		MockInjection()
	default:
		task.SetAuthTools(auth.AdminAuthService)
	}
}

func MockInjection() {
	task.SetAuthTools(auth.MockAuthService)
}
