package task

type DbInjectionRules struct {
	ID          uint64 `json:"id" gorm:"column:id"`
	RuleName    string `json:"rule_name" gorm:"column:rule_name"`
	RuleContent string `json:"rule_content" gorm:"column:rule_content"`
	Switch      int    `json:"switch" gorm:"column:switch"`
	Ct          int64  `json:"ct" gorm:"column:ct"`
	Ut          int64  `json:"ut" gorm:"column:ut"`
}

type DbInjectionSubtask struct {
	ID          int64  `json:"id" gorm:"column:id"`
	TaskID      int64  `json:"task_id" gorm:"column:task_id"`
	TaskType    string `json:"task_type" gorm:"column:task_type"`
	DbName      string `json:"db_name" gorm:"column:db_name"`
	ClusterName string `json:"cluster_name" gorm:"column:cluster_name"`
	Ct          int64  `json:"ct" gorm:"column:ct"`
	Ut          int64  `json:"ut" gorm:"column:ut"`

	ExecItems []DbInjectionExecItem `json:"exec_items" gorm:"-"`
}

type DbInjectionExecItem struct {
	ID             int64  `json:"id" gorm:"column:id"`
	TaskID         int64  `json:"task_id" gorm:"column:task_id"`
	SubTaskID      int64  `json:"sub_task_id" gorm:"column:sub_task_id"`
	SQLContent     string `json:"sql_content" gorm:"column:sql_content"`
	Remark         string `json:"remark" gorm:"column:remark"`
	AffectRows     int    `json:"affect_rows" gorm:"column:affect_rows"`
	InjectComments string `json:"inject_comments" gorm:"column:inject_comments"`
	Status         string `json:"status" gorm:"column:status"`
	Ct             int64  `json:"ct" gorm:"column:ct"`
	Ut             int64  `json:"ut" gorm:"column:ut"`
	Et             int64  `json:"et" gorm:"column:et"`
	ExecInfo       string `json:"exec_info" gorm:"column:exec_info"`
	BackupInfo     string `json:"backup_info" gorm:"column:exec_info"`
	BackupStatus   string `json:"backup_status" gorm:"column:backup_status"`
	BackupID       int64  `json:"backup_id" gorm:"column:backup_id"`
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

	//中止状态
	Reject     Status = "reject"
	Cancel            = "cancel"
	ExecCancel        = "execCancel"

	//子项状态
	SkipExec Status = "skipExec"
)

type ItemStatus = string

const (
	ItemFailed       ItemStatus = "failed"
	ItemBackupFailed            = "backupFailed"
	ItemBackupSuccess            = "backupSuccess"
	ItemSuccess                 = "success"
	ItemSkipped                 = "skipped"
)

type ExecType = string

const (
	DML ExecType = "DML"
	DDL          = "DDL"
)

type Action = string

const (
	EditItem Action = "editItem"
	DoCancel        = "cancel"
	SkipAt          = "skipAt"
	BeginAt         = "beginAt"
	Progress        = "progress"
	DoReject        = "reject"
)
