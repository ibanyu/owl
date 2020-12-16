package task

type DbInjectionSubtask struct {
	ID          int64  `json:"id" gorm:"column:id"`
	TaskID      int64  `json:"task_id" gorm:"column:task_id"`
	TaskType    string `json:"task_type" gorm:"column:task_type"`
	DbName      string `json:"db_name" gorm:"column:db_name"`
	ClusterName string `json:"cluster_name" gorm:"column:cluster_name"`

	ExecItems []DbInjectionExecItem `json:"exec_items" gorm:"-"`
}

type DbInjectionExecItem struct {
	ID           int64  `json:"id" gorm:"column:id"`
	TaskID       int64  `json:"task_id" gorm:"column:task_id"`
	SubtaskID    int64  `json:"subtask_id" gorm:"column:subtask_id"`
	SQLContent   string `json:"sql_content" gorm:"column:sql_content"`
	Remark       string `json:"remark" gorm:"column:remark"`
	AffectRows   int    `json:"affect_rows" gorm:"column:affect_rows"`
	RuleComments string `json:"rule_comments" gorm:"column:rule_comments"`
	Status       string `json:"status" gorm:"column:status"`
	ExecInfo     string `json:"exec_info" gorm:"column:exec_info"`
	BackupInfo   string `json:"backup_info" gorm:"column:backup_info"`
	BackupStatus string `json:"backup_status" gorm:"column:backup_status"`
	BackupID     int64  `json:"backup_id" gorm:"column:backup_id"`

	Ut int64 `json:"ut" gorm:"column:ut"`
	Et int64 `json:"et" gorm:"column:et"`

	DBName      string `json:"db_name" gorm:"-"`
	ClusterName string `json:"cluster_name" gorm:"-"`
	TaskType    string `json:"task_type" gorm:"-"`
}

type Status = string

const (
	//顺序递进
	CheckFailed Status = "checkFailed"
	CheckPass          = "checkPass"
	ReviewPass         = "reviewPass"
	DBAPass            = "dbaPass" // exec wait
	ExecFailed         = "execFailed"
	ExecSuccess        = "execSuccess"

	//终止状态
	Reject     Status = "reject"
	Cancel            = "cancel"
	ExecCancel        = "execCancel"

	//子项状态
	SkipExec Status = "skipExec"
)

type ItemStatus = string

const (
	ItemFailed  ItemStatus = "failed"
	ItemSuccess            = "success"
	ItemSkipped            = "skipped"

	ItemNoBackup              ItemStatus = "rollBackFailed"
	ItemBackupSuccess                    = "backupSuccess"
	ItemBackupFailed                     = "backupFailed"
	ItemRollBackFailed                   = "rollBackFailed"
	ItemRollBackSuccess                  = "rollBackSuccess"
	ItemAlreadyRollBack                  = "rollBackSuccess"
	ItemAlreadyRollBackFailed            = "rollBackSuccess"
)

type Type = string

const (
	DML Type = "DML"
	DDL      = "DDL"
)

type Action = string

const (
	EditItem Action = "editItem"
	DelItem         = "delItem"
	DoCancel        = "cancel"
	SkipAt          = "skipAt"
	BeginAt         = "beginAt"
	Progress        = "progress"
	DoReject        = "reject"
)

//todo, list not in, history in;
//列表及历史 添加能见过滤；创建者，reviewer
//本地部署，添加dockerfile；
func HistoryStatus() []ItemStatus {
	return []ItemStatus{Reject, Cancel, ExecCancel, ExecFailed, ExecSuccess}
}
