package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/checker"
)

func LisRule(ctx *gin.Context) Resp {
	rules := checker.ListRules()
	return Resp{Data: ListData{
		Items:  rules,
		Total:  len(rules),
		More:   false,
		Offset: 0,
	}}
}

func UpdateRuleStatus(ctx *gin.Context) Resp {
	f := "UpdateRuleStatus()-->"

	taskName := ctx.Query("task_name")
	action := ctx.Query("action")
	if taskName == "" || action == "" {
		return Resp{Message: fmt.Sprintf("%s, get param failed ", f), Code: code.ParamInvalid}
	}

	if err := checker.UpdateRuleStatus(taskName, action); err != nil {
		return Resp{Message: fmt.Sprintf("%s,update RuleStatus failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{}
}
