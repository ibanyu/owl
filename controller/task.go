package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

func ListTask(ctx *gin.Context) Resp {
	f := "ListTask() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	task, count, err := task.ListTask(&page)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: list Cabinet err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: ListData{
		Items: task,
		Count: count,
	}}
}

func GetTask(ctx *gin.Context) Resp {
	f := "GetTask() -->"
	user, err := ctx.Cookie("user")
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: get user failed, err: %s", f, err.Error()), Code: code.UserNotFound}
	}

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Resp{
			Message: fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr),
			Code:    code.ParamInvalid,
		}
	}

	task, err := task.GetTask(id, user)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: get task failed, err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: task}
}

func UpdateTask(ctx *gin.Context) Resp {
	f := "UpdateTask()-->"
	var taskParam task.DbInjectionTask
	if err := ctx.BindJSON(&taskParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}

	if err := task.UpdateTask(&taskParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, update task failed :%s", f, err.Error()), Code: code.InternalErr}
	}
	return Resp{}
}

func AddTask(ctx *gin.Context) Resp {
	f := "AddTask()-->"
	var taskParam task.DbInjectionTask
	if err := ctx.BindJSON(&taskParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}

	id, err := task.AddTask(&taskParam)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s, get task failed :%s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: id}
}
