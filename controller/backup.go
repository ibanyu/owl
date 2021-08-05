package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

func ListRollbackData(ctx *gin.Context) Resp {
	f := "GetTask() -->"

	var req task.RollBackReq
	if err := ctx.BindJSON(&req); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	rollBackData, err := task.ListRollbackData(&req)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: get rollBackData failed, err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: rollBackData}
}

func Rollback(ctx *gin.Context) Resp {
	f := "Rollback()-->"
	var req task.RollBackReq
	if err := ctx.BindJSON(&req); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	req.Executor = ctx.MustGet("user").(string)
	if err := task.Rollback(&req); err != nil {
		return Resp{Message: fmt.Sprintf("%s: rollback failed, err: %s", f, err.Error()), Code: code.InternalErr}
	}
	return Resp{}
}
