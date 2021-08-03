package task

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

func Exec(paramTask, dbTask *DbInjectionTask) error {
	if paramTask.Et == 0 {
		return ExecTask(paramTask, dbTask)
	}

	err := taskDao.UpdateTask(&DbInjectionTask{
		ID:       dbTask.ID,
		Status:   ExecWait,
		Ut:       time.Now().Unix(),
		Et:       paramTask.Et,
		Executor: paramTask.Executor,
	})
	if err != nil {
		return fmt.Errorf("before exec a cron task, persist it err:%s", err.Error())
	}

	return nil
}

//exec and update status
//exec from head, skip at some one, or begin at some one
func ExecTask(paramTask, dbTask *DbInjectionTask) error {
	//todo, 注意检查获取，执行的顺序

	startId, err := getExecStartId(paramTask.Action, dbTask.ExecItems, paramTask.ExecItem)
	if err != nil {
		return err
	}

	// mean need't exec task
	if startId < 0 {
		return taskDao.UpdateTask(&DbInjectionTask{
			ID:     paramTask.ID,
			Status: ExecSuccess,
			Ut:     time.Now().Unix(),
		})
	}

	//exec task
	jump := true
	failed := false
	beginTime := time.Now().Unix()
	for _, subTask := range dbTask.SubTasks { //todo
		dbInfo, err := dbTool.GetDBConn(subTask.DbName, subTask.ClusterName)
		if err != nil {
			return err
		}

		for _, item := range subTask.ExecItems {
			if item.ID != startId && jump {
				continue
			}
			jump = false

			err := BackupAndExec(dbInfo.DB, &item, subTask.TaskType)
			if err != nil {
				failed = true
				err := taskDao.UpdateTask(&DbInjectionTask{
					ID:       dbTask.ID,
					Status:   ExecFailed,
					Et:       beginTime,
					Executor: paramTask.Executor,
					ExecInfo: err.Error(),
				})
				if err != nil {
					logger.Errorf("after exec failed, update task status to failed err, err： %s", err.Error())
				}
				break
			}
		}

		dbInfo.CloseConn()
	}

	if !failed {
		err = taskDao.UpdateTask(&DbInjectionTask{
			ID:       dbTask.ID,
			Status:   ExecSuccess,
			Et:       beginTime,
			Ft:       time.Now().Unix(),
			Executor: paramTask.Executor,
		})
		if err != nil {
			logger.Errorf("after exec, update task status to success err, err： %s", err.Error())
		}
	}

	return nil
}

// backup, exec, update status
func BackupAndExec(db *sql.DB, item *DbInjectionExecItem, taskType string) error {
	execBackup, backupId, backupErr := backup(db, taskType, item.SQLContent)
	if !execBackup {
		item.BackupStatus = ItemNoBackup
	} else if backupErr == nil {
		item.BackupStatus = ItemBackupSuccess
	} else {
		err := subTaskDao.UpdateItem(&DbInjectionExecItem{
			ID:           item.ID,
			BackupStatus: ItemBackupFailed,
			BackupInfo:   backupErr.Error(),
		})
		if err != nil {
			logger.Errorf("while backup failed, update item backup status err, %s", err.Error())
		}

		if !config.Conf.Server.ExecNoBackup {
			return fmt.Errorf("backup err: %s", backupErr.Error())
		}
	}

	result, err := db.Exec(item.SQLContent)
	if err != nil {
		item.Status = ItemFailed
		item.ExecInfo = err.Error()
	} else {
		item.Status = ItemSuccess
		item.BackupID = backupId
		item.ExecInfo = fmtExecInfo(result)
	}

	item.Et = time.Now().Unix()
	updateStatusErr := subTaskDao.UpdateItem(item)
	if updateStatusErr != nil {
		logger.Errorf("after exec, update execItem status err, err： %s", updateStatusErr.Error())
	}
	return err
}

func fmtExecInfo(result sql.Result) string {
	if result == nil {
		return ""
	}

	affect, _ := result.RowsAffected()
	lastId, _ := result.LastInsertId()
	return fmt.Sprintf("affect nums: %d, last insert id: %d", affect, lastId)
}

func getExecStartId(action Action, subItems []DbInjectionExecItem, targetItem *DbInjectionExecItem) (int64, error) {
	switch action {
	case Progress:
		for _, v := range subItems {
			if v.Status != ItemSuccess {
				return v.ID, nil
			}
		}
		return -1, nil
	case BeginAt:
		return targetItem.ID, nil
	case SkipAt:
		find := false
		for _, v := range subItems {
			if find {
				return v.ID, nil
			}
			if v.ID == targetItem.ID {
				find = true
				err := subTaskDao.UpdateItem(&DbInjectionExecItem{ID: v.ID, Status: ItemSkipped})
				if err != nil {
					logger.Errorf("update task status to skip failed, taskId: %d", v.ID)
				}
			}
		}

		//跳过的是最后一个，则不执行
		if find {
			return -1, nil
		} else {
			return 0, fmt.Errorf("execute skip task, target not found, targeId: %d", targetItem.ID)
		}
	default:
		return 0, fmt.Errorf("execute task err, type not found, type: %d", action)
	}
}

//todo
func fmtExecItemFromOneTask(task *DbInjectionTask) (items []DbInjectionExecItem) {
	for _, subTask := range task.SubTasks {
		for _, v := range subTask.ExecItems {
			items = append(items, v)
		}
	}

	return
}
