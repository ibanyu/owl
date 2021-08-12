package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ibanyu/owl/code"
	"github.com/ibanyu/owl/service/checker"
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

	params := struct {
		Name   string `json:"name" binding:"required"`
		Action string `json:"action" binding:"required"`
	}{}
	if err := ctx.BindJSON(&params); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	if err := checker.UpdateRuleStatus(params.Name, params.Action); err != nil {
		return Resp{Message: fmt.Sprintf("%s,update RuleStatus failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{}
}
