package controller

import (
	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
)

func LisRule(ctx *gin.Context) Resp {
	rules := checker.ListRules()
	return Resp{Data: ListData{
		Items: rules,
		Count: len(rules),
	}}
}
