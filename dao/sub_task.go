package dao

import "gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"

type SubTaskDaoImpl struct {
}

var SubTask SubTaskDaoImpl

func (SubTaskDaoImpl) UpdateItem(item *task.DbInjectionExecItem) error {
	return GetDB().Where("id = ?", item.ID).Update(item).Error
}

func (SubTaskDaoImpl) UpdateItemByBackupId(item *task.DbInjectionExecItem) error {
	return GetDB().Where("backup_id = ?", item.BackupID).Update(item).Error
}
