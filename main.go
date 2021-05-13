package main

import (
	"flag"
	"log"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/dao"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/injection"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/router"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

var Version string

func main() {
	flag.Parse()
	log.Println("version:", Version)
	config.InitConfig("")
	logger.InitLog(config.Conf.Server.LogDir, "dbinjection.log", config.Conf.Server.LogLevel)
	dao.InitDB()
	injection.Injection()
	checker.InitRuleStatus()

	go controller.ExecWaitTask()
	router.Run()
}
