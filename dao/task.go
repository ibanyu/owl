package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
)

type TaskDaoImpl struct {
}

var Task TaskDaoImpl

func (TaskDaoImpl) AddTask(task *task.DbInjectionTask) (int64, error) {
	tx := GetDB().Begin()
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, subTask := range task.SubTasks {
		subTask.TaskID = task.ID
		if err := tx.Create(&subTask).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, item := range subTask.ExecItems {
			item.SubtaskID = subTask.ID
			item.TaskID = task.ID
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	return task.ID, tx.Commit().Error
}

func (TaskDaoImpl) UpdateTask(task *task.DbInjectionTask) error {
	return GetDB().Model(task).Where("id = ?", task.ID).Update(task).Error
}

func (TaskDaoImpl) ListTask(page *service.Pagination, isDBA bool) ([]task.DbInjectionTask, int, error) {
	condition := "(name like ? or creator like ?) and status not in (?)"
	if !isDBA {
		condition = fmt.Sprintf("(creator = %s or reviewer = %s) and ", page.Operator, page.Operator) + condition
	}

	page.Key = "%" + page.Key + "%"
	var count int
	if err := GetDB().Model(&task.DbInjectionTask{}).Where(condition,
		page.Key, page.Key, task.HistoryStatus()).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var tasks []task.DbInjectionTask
	if err := GetDB().Order("ct desc").Offset(page.Offset).Limit(page.Limit).
		Find(&tasks, condition, page.Key, page.Key, task.HistoryStatus()).Error; err != nil {
		return nil, 0, err
	}

	for idx, taskV := range tasks {
		formattedItems, _, err := getTaskExecItems(GetDB(), &taskV)
		if err != nil {
			return nil, 0, err
		}

		tasks[idx].ExecItems = formattedItems
	}

	return tasks, count, nil
}

func (TaskDaoImpl) ListHistoryTask(page *service.Pagination, isDBA bool) ([]task.DbInjectionTask, int, error) {
	condition := "(name like ? or creator like ?) and status in (?)"
	if !isDBA {
		condition = fmt.Sprintf(" and (creator = %s or reviewer = %s)", page.Operator, page.Operator) + condition
	}

	page.Key = "%" + page.Key + "%"
	var count int
	if err := GetDB().Model(&task.DbInjectionTask{}).Where(condition,
		page.Key, page.Key, task.HistoryStatus()).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var tasks []task.DbInjectionTask
	if err := GetDB().Order("ct desc").Offset(page.Offset).Limit(page.Limit).
		Find(&tasks, condition, page.Key, page.Key, task.HistoryStatus()).Error; err != nil {
		return nil, 0, err
	}

	for idx, taskV := range tasks {
		formattedItems, _, err := getTaskExecItems(GetDB(), &taskV)
		if err != nil {
			return nil, 0, err
		}

		tasks[idx].ExecItems = formattedItems
	}

	return tasks, count, nil
}

func getTaskExecItems(db *gorm.DB, taskP *task.DbInjectionTask) ([]task.DbInjectionExecItem, []task.DbInjectionSubtask, error) {
	var formattedItems []task.DbInjectionExecItem
	var subTasks []task.DbInjectionSubtask
	if err := db.Find(&subTasks, "task_id = ?", taskP.ID).Error; err != nil {
		return nil, nil, err
	}

	for idx, subTask := range subTasks {
		var items []task.DbInjectionExecItem
		if err := db.Find(&items, "subtask_id = ?", subTask.ID).Error; err != nil {
			return nil, nil, err
		}
		subTasks[idx].ExecItems = items

		for _, v := range items {
			v.DBName = subTask.DbName
			v.ClusterName = subTask.ClusterName
			v.TaskType = subTask.TaskType
			formattedItems = append(formattedItems, v)
		}
	}

	return formattedItems, subTasks, nil
}

func (TaskDaoImpl) GetTask(id int64) (*task.DbInjectionTask, error) {
	var task task.DbInjectionTask
	if err := GetDB().First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	formattedItems, subTasks, err := getTaskExecItems(GetDB(), &task)
	if err != nil {
		return nil, err
	}

	task.SubTasks = subTasks
	task.ExecItems = formattedItems
	return &task, nil
}
