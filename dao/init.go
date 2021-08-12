package dao

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/ibanyu/owl/config"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func InitDB() {
	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4", config.Conf.DB.User,
		config.Conf.DB.Password, config.Conf.DB.Address, config.Conf.DB.Port, config.Conf.DB.DBName)

	var err error
	DB, err = gorm.Open("mysql", conn)
	if err != nil {
		log.Fatal("init DB error : ", err.Error())
	}

	if config.Conf.DB.MaxIdleConn > 0 {
		DB.DB().SetMaxIdleConns(config.Conf.DB.MaxIdleConn)
	}
	if config.Conf.DB.MaxOpenConn > 0 {
		DB.DB().SetMaxOpenConns(config.Conf.DB.MaxOpenConn)
	}
	DB.SingularTable(true)
	DB.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {})

	if config.Conf.Server.ShowSql {
		DB.LogMode(true)
	}
}
