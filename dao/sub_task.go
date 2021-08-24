package dao

import "github.com/ibanyu/owl/service/task"

type SubTaskDaoImpl struct {
}

var SubTask SubTaskDaoImpl

func (SubTaskDaoImpl) UpdateItem(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("id = ?", item.ID).Update(item).Error
}

func (SubTaskDaoImpl) DelItem(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("id = ?", item.ID).Delete(item).Error
}

func (SubTaskDaoImpl) UpdateItemByBackupId(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("backup_id = ?", item.BackupID).Update(item).Error
}
