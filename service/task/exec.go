package task

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

//exec and update status
//exec from head, skip at some one, or begin at some one
func ExecTask(paramTask *DbInjectionTask) error {
	//todo, 注意检查获取，执行的顺序
	task, err := taskDao.GetTask(paramTask.ID)
	if err != nil {
		return err
	}

	task.ExecItems = fmtExecItemFromOneTask(task)

	startId, err := getExecStartId(paramTask.Action, task.ExecItems, &paramTask.ExecItem)
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
	for _, subTask := range task.SubTasks {
		dbConn, err := dbTool.GetDBConn(subTask.DbName, subTask.ClusterName)
		if err != nil {
			return err
		}

		for _, item := range subTask.ExecItems {
			if item.ID != startId && jump {
				continue
			}
			jump = false

			err := BackupAndExec(dbConn, &item, subTask.TaskType)
			if err != nil {
				err = taskDao.UpdateTask(&DbInjectionTask{
					ID:       task.ID,
					Status:   ExecFailed,
					Et:       beginTime,
					Executor: paramTask.Executor,
					ExecInfo: err.Error(),
				})
				if err != nil {
					logger.Errorf("after exec, update task status err, err： %s", err.Error())
				}
				break
			}
		}
	}

	if !failed {
		err = taskDao.UpdateTask(&DbInjectionTask{
			ID:       task.ID,
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
	backSuccess, backupId, backupErr := backup(db, taskType, item.SQLContent)
	if backSuccess {
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
		item.ExecInfo = fmt.Sprintf("%v", result)
	}

	item.Et = time.Now().Unix()
	updateStatusErr := subTaskDao.UpdateItem(item)
	if updateStatusErr != nil {
		logger.Errorf("after exec, update execItem status err, err： %s", updateStatusErr.Error())
	}
	return err
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

func fmtExecItemFromOneTask(task *DbInjectionTask) (items []DbInjectionExecItem) {
	for _, subTask := range task.SubTasks {
		for _, v := range subTask.ExecItems {
			items = append(items, v)
		}
	}

	return
}
