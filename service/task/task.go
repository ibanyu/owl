package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/sql_util"
)

type DbInjectionTask struct {
	ID            int64  `json:"id" gorm:"column:id"`
	Name          string `json:"name" gorm:"column:name"`
	Status        string `json:"status" gorm:"column:status"`
	Creator       string `json:"creator" gorm:"column:creator"`
	Reviewer      string `json:"reviewer" gorm:"column:reviewer"`
	Executor      string `json:"executor" gorm:"column:executor"`
	ExecInfo      string `json:"exec_info" gorm:"column:exec_info"`
	RejectContent string `json:"reject_content" gorm:"column:reject_content"`
	Ct            int64  `json:"ct" gorm:"column:ct"`
	Ut            int64  `json:"ut" gorm:"column:ut"`
	Et            int64  `json:"et" gorm:"column:et"`
	Ft            int64  `json:"ft" gorm:"column:ft"`

	SubTasks  []DbInjectionSubtask  `json:"sub_tasks" gorm:"-"`
	ExecItems []DbInjectionExecItem `json:"exec_items" gorm:"-"`
	ExecItem  *DbInjectionExecItem  `json:"exec_item" gorm:"-"`
	EditAuth  *EditAuth             `json:"edit_auth" gorm:"-"`

	Action Action `json:"action" gorm:"-"`
}

type TaskDao interface {
	AddTask(task *DbInjectionTask) (int64, error)
	UpdateTask(task *DbInjectionTask) error
	ListTask(pagination *service.Pagination) ([]DbInjectionTask, int, error)
	GetTask(id int64) (*DbInjectionTask, error)
}

var taskDao TaskDao

func SetTaskDao(impl TaskDao) {
	taskDao = impl
}

type SubTaskDao interface {
	UpdateItem(item *DbInjectionExecItem) error
	UpdateItemByBackupId(item *DbInjectionExecItem) error
}

var subTaskDao SubTaskDao

func SetSubTaskDao(impl SubTaskDao) {
	subTaskDao = impl
}

func AddTask(task *DbInjectionTask) (int64, error) {
	reviewer, err := authTool.GetReviewer(task.Creator)
	if err != nil {
		return 0, err
	}
	task.Reviewer = reviewer

	// 拆分
	for idx, subTask := range task.SubTasks {
		newSubTask, err := splitSubTaskExecItems(&subTask)
		if err != nil {
			return 0, err
		}
		task.SubTasks[idx] = *newSubTask
	}

	if err := checkExecItemNum(task); err != nil {
		return 0, err
	}

	// check sql
	checkPass := true
	for idx, subTask := range task.SubTasks {
		dbInfo, err := dbTool.GetDBConn(subTask.DbName, subTask.ClusterName)
		if err != nil {
			return 0, err
		}

		for itemIdx, item := range subTask.ExecItems {
			pass, suggestion, affectRow := checker.SqlCheck(item.SQLContent, "", "", dbInfo)
			if affectRow > 0 {
				task.SubTasks[idx].ExecItems[itemIdx].AffectRows = 0
			}
			if !pass {
				checkPass = false
				task.SubTasks[idx].ExecItems[itemIdx].RuleComments = suggestion
			}
		}

		dbInfo.CloseConn()
	}

	if checkPass {
		task.Status = CheckPass
	}
	task.Ct = time.Now().Unix()

	// add task
	return taskDao.AddTask(task)
}

func splitSubTaskExecItems(subTask *DbInjectionSubtask) (*DbInjectionSubtask, error) {
	var items []DbInjectionExecItem
	for _, execItem := range subTask.ExecItems {
		sqls, err := sql_util.SplitMultiSql(execItem.SQLContent)
		if err != nil {
			return nil, err
		}
		for _, v := range sqls {
			newItem := execItem
			newItem.SQLContent = v
			items = append(items, newItem)
		}
	}
	subTask.ExecItems = items
	return subTask, nil
}

func checkExecItemNum(task *DbInjectionTask) error {
	num := 0
	for _, v := range task.SubTasks {
		for range v.ExecItems {
			num++
		}
	}
	if num > config.Conf.Server.NumOnceLimit {
		return fmt.Errorf("exec too many sql once, sqlNum: %d, limit: %d", num, config.Conf.Server.NumOnceLimit)
	}
	return nil
}

func UpdateTask(task *DbInjectionTask) error {
	dbTask, err := taskDao.GetTask(task.ID)
	if err != nil {
		return err
	}

	isReviewer := strings.Contains(dbTask.Reviewer, task.Executor)
	isDba, err := authTool.IsDba(task.Executor)
	if err != nil {
		return err
	}

	//对于执行变更，检查权限
	if task.Action == BeginAt || task.Action == SkipAt || (dbTask.Status == DBAPass && task.Action == Progress) {
		if !(isReviewer && allIsDmlTask(task)) && !isDba {
			return errors.New("auth invalid")
		}
	}

	switch task.Action {
	case EditItem:
		return subTaskDao.UpdateItem(task.ExecItem)
	case DoCancel:
		return doCancel(task, dbTask, isDba)
	case SkipAt, BeginAt:
		return ExecTask(task, dbTask)
	case Progress:
		return ProgressEdit(task, dbTask)
	case DoReject:
		return taskDao.UpdateTask(&DbInjectionTask{
			ID:            task.ID,
			Status:        Reject,
			Executor:      task.Executor,
			RejectContent: task.RejectContent,
			Ut:            time.Now().Unix(),
		})
	default:
		return fmt.Errorf("action type not found, action: %s", task.Action)
	}
}

func doCancel(task, dbTask *DbInjectionTask, isDba bool) error {
	switch {
	case dbTask.Creator == task.Executor:
		task.Status = Cancel
	case isDba:
		task.Status = ExecCancel
	default:
		return errors.New("no auth to do cancel")
	}

	return taskDao.UpdateTask(&DbInjectionTask{
		ID:       task.ID,
		Status:   task.Status,
		Executor: task.Executor,
	})
}

func ProgressEdit(task, dbTask *DbInjectionTask) error {
	switch dbTask.Status {
	case CheckPass:
		task.Status = ReviewPass
	case ReviewPass:
		task.Status = DBAPass
	case DBAPass:
		return ExecTask(task, dbTask)
	default:
		return fmt.Errorf("progress failed, task status invalid, status: %s", dbTask.Status)
	}

	return taskDao.UpdateTask(&DbInjectionTask{
		ID:       task.ID,
		Status:   task.Status,
		Ut:       time.Now().Unix(),
		Executor: task.Executor,
	})
}

func ListTask(pagination *service.Pagination) ([]DbInjectionTask, int, error) {
	tasks, count, err := taskDao.ListTask(pagination)
	if err != nil {
		return nil, 0, err
	}

	for i, v := range tasks {
		isDba, err := authTool.IsDba(pagination.Operator)
		if err != nil {
			return nil, 0, err
		}
		tasks[i].EditAuth = GetTaskOperateAuth(false, v.Creator == pagination.Operator, strings.Contains(v.Reviewer, pagination.Operator), isDba, &v)
	}

	return tasks, count, nil
}

func GetTask(id int64, operator string) (*DbInjectionTask, error) {
	isDba, err := authTool.IsDba(operator)
	if err != nil {
		return nil, err
	}

	task, err := taskDao.GetTask(id)
	if err != nil {
		return nil, err
	}
	task.SubTasks = nil

	//task.ExecItems = fmtExecItemFromOneTask(task)
	task.EditAuth = GetTaskOperateAuth(true, operator == task.Creator, isDba, strings.Contains(task.Reviewer, operator), task)
	return task, nil
}
