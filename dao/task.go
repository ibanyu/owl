package dao

import (
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
		if err := tx.Create(subTask).Error; err != nil {
			tx.Rollback()
			return 0, err

			for _, item := range subTask.ExecItems {
				item.SubtaskID = subTask.ID
				item.TaskID = task.ID
				if err := tx.Create(item).Error; err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	return task.ID, tx.Commit().Error
}

func (TaskDaoImpl) UpdateTask(task *task.DbInjectionTask) error {
	return GetDB().Where("id = ?", task.ID).Update(task).Error
}

const listTaskCondition = "name like ? or `type` = ?"

func (TaskDaoImpl) ListTask(page *service.Pagination) ([]task.DbInjectionTask, int, error) {
	page.Key = "%" + page.Key + "%"

	var count int
	if err := GetDB().Model(&task.DbInjectionTask{}).Where(listTaskCondition,
		page.Key, page.Key).Count(&count).Error;
		err != nil {
		return nil, 0, err
	}

	var modules []task.DbInjectionTask
	return modules, count, GetDB().Order("ct desc").Offset(page.Offset).Limit(page.Limit).Find(&modules, listTaskCondition, page.Key, page.Key).Error
}

func (TaskDaoImpl) GetTask(id int64) (*task.DbInjectionTask, error) {
	var task task.DbInjectionTask
	return &task, GetDB().First(&task, "id = ?", id).Error
}
