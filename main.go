package main

import (
	"flag"
	"log"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/dao"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/injection"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/router"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

var Version string

func main() {
	flag.Parse()
	log.Println("version:", Version)
	logger.InitLog()
	config.InitConfig("")
	dao.InitDB()
	injection.Injection()

	router.Run()
}
