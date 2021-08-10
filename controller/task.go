package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

func ListExecTask(ctx *gin.Context) Resp {
	f := "ListExecTask() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	page.Operator = ctx.MustGet("user").(string)
	task, count, err := task.ListTask(&page, task.ExecStatus())
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: list ListExecTask err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: ListData{
		Items:  setTypeNameFroTasks(task),
		Total:  count,
		More:   count > page.Offset+page.Limit,
		Offset: page.Offset,
	}}
}

func ListReviewerTask(ctx *gin.Context) Resp {
	f := "ListExecTask() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	page.Operator = ctx.MustGet("user").(string)
	task, count, err := task.ListTask(&page, task.ReviewerStatus())
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: list ListExecTask err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: ListData{
		Items:  setTypeNameFroTasks(task),
		Total:  count,
		More:   count > page.Offset+page.Limit,
		Offset: page.Offset,
	}}
}

func ListHistoryTask(ctx *gin.Context) Resp {
	f := "ListHistoryTask() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	page.Operator = ctx.MustGet("user").(string)
	task, count, err := task.ListHistoryTask(&page)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: list ListHistoryTask err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: ListData{
		Items:  setTypeNameFroTasks(task),
		Total:  count,
		More:   count > page.Offset+page.Limit,
		Offset: page.Offset,
	}}
}

func GetTask(ctx *gin.Context) Resp {
	f := "GetTask() -->"

	user := ctx.MustGet("user").(string)
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

	return Resp{Data: setTypeNameFroTask(task)}
}

func UpdateTask(ctx *gin.Context) Resp {
	f := "UpdateTask()-->"
	var taskParam task.DbInjectionTask
	if err := ctx.BindJSON(&taskParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}

	taskParam.Executor = ctx.MustGet("user").(string)
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

	if err := task.CheckTaskType(&taskParam); err != nil {
		return Resp{Message: err.Error(), Code: code.ParamInvalid}
	}

	taskParam.Creator = ctx.MustGet("user").(string)
	id, err := task.AddTask(&taskParam)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s, get task failed :%s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: struct {
		Id int64 `json:"id"`
	}{id}}
}

func ExecWaitTask() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("The execWaitTask goroutine panic, err:%s", err)
			// keep the goroutine running
			ExecWaitTask()
		}
	}()

	for {
		<-time.After(time.Minute)

		waitTasks, _, err := task.GetExecWaitTask()
		if err != nil {
			logger.Errorf("the goroutine get exec wait tasks err:%v", err)
		}
		for _, waitTask := range waitTasks {
			countDown := waitTask.Et - time.Now().Unix()
			if countDown <= 0 {
				waitTask.Action = task.Progress
				if err := task.ExecTask(&waitTask, &waitTask); err != nil{
					logger.Errorf("while exec task in cron err: %s", err.Error())
				}
			}
		}
	}
}

func setTypeNameFroTasks(tasks []task.DbInjectionTask) []task.DbInjectionTask {
	for i, _ := range tasks {
		setTypeNameFroTask(&tasks[i])
	}
	return tasks
}

func setTypeNameFroTask(oneTask *task.DbInjectionTask) *task.DbInjectionTask {
	for i, v := range oneTask.ExecItems {
		oneTask.ExecItems[i].TaskType = task.GetTypeName(v.TaskType)
	}
	return oneTask
}
