package main

import (
	"flag"
	"log"

	"github.com/ibanyu/owl/config"
	"github.com/ibanyu/owl/controller"
	"github.com/ibanyu/owl/dao"
	"github.com/ibanyu/owl/injection"
	"github.com/ibanyu/owl/router"
	"github.com/ibanyu/owl/service/checker"
	"github.com/ibanyu/owl/util/logger"
)

var Version string

func main() {
	flag.Parse()
	log.Println("version:", Version)
	config.InitConfig("")
	logger.InitLog(config.Conf.Server.LogDir, "owl.log", config.Conf.Server.LogLevel)
	dao.InitDB()
	injection.Injection()
	checker.InitRuleStatus()

	go controller.ExecWaitTask()
	router.Run()
}
