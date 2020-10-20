package main

import (
	"log"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/dao"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/injection"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/router"
)

var Version string

func main() {
	log.Println("version:", Version)
	config.InitConfig("")
	dao.InitDB()
	injection.Injection()

	router.Run()
}
