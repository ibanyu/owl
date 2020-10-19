package task

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/pingcap/parser/ast"
	"github.com/shawnfeng/sutil/slog"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/sql_util"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xsql/scanner"
)

type BackupDao interface {
	AddBackup(backup *DbInjectionBackup) (int64, error)
	UpdateBackupStatus(backup *DbInjectionBackup) error
	GetBackupInfoById(id int64) (*DbInjectionBackup, error)
}

var backupDao BackupDao

func SetBackupDao(impl BackupDao) {
	backupDao = impl
}

type DbInjectionBackup struct {
	ID           int64  `json:"id" bdb:"id"`
	Data         string `json:"data" bdb:"data"`
	Ct           int64  `json:"ct" bdb:"ct"`
	RollbackTime uint64 `json:"rollback_time" bdb:"rollback_time"`
	RollbackUser string `json:"rollback_user" bdb:"rollback_user"`
	IsRollBacked int    `json:"is_roll_backed" bdb:"is_roll_backed"`
}

const (
	splitFieldAsterisk = "*"
	splitRowNumberSign = "#"
	NUL                = ""
)
const (
	replaceAsterisk   = "0d65760d0deb56cb59a49ebbe0939cf3"
	replaceNumberSign = "c13c1fce3ca97e765860002118ca3bb8"
	replaceNUL        = "4252a939eaac0e4f457f296e60882ebe"
)

func backup(db *sql.DB, taskType, sql string) (execBackup bool, backupId int64, err error) {
	if taskType != DML ||
		!needBackup(sql) {
		slog.Infof("not a update or delete operate , don't backup")
		return
	}
	execBackup = true
	slog.Infof("start backup data ...")

	selectSql, tableName, err := getSqlInfo(sql)
	if err != nil {
		slog.Errorf("convert sql to select err: %s", err.Error())
		return
	}
	isEmpty, backupId, err := fetchAndStoreBackupInfo(db, selectSql, tableName)
	if isEmpty {
		err = errors.New("backup nothing, condition match nothing")
	}
	return
}

func needBackup(sql string) bool {
	stmtNodes, _, _ := sql_util.Parser.Parse(sql, "", "")

	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			return true
		default:
			return false
		}
	}
	return false
}

func fetchAndStoreBackupInfo(db *sql.DB, selectSql, tableName string) (isEmpty bool, backupId int64, err error) {
	rows, err := db.Query(selectSql)
	if err != nil {
		slog.Warnf("exec backup db exec err:%v", err)
		return
	}

	column, err := sql_util.GetTableColumn(tableName, db)
	if err != nil {
		return
	}

	dataStr := formatData(rows, column)
	if strings.TrimSpace(dataStr) == "" {
		isEmpty = true
		slog.Warnf("while backup, condition match nothing")
		return
	}

	backupId, err = backupDao.AddBackup(&DbInjectionBackup{
		Data: dataStr,
		Ct:   time.Now().Unix(),
	})
	return
}

func getSqlInfo(sql string) (selectSql string, tableName string, err error) {
	tableName, err = sql_util.GetTableName(sql)

	where := sql_util.GetSqlAfterWhere(sql)
	selectSql = fmt.Sprintf("select * from %s where %s", tableName, sql_util.HandelKeyWorldForCondition(where))
	slog.Infof("build backup select sql : %s ", selectSql)

	return
}

//map 乱序, 需要按列顺序存
// 用#分割行,用*分割字段
// 需要把字段中的 # * 空字符替换掉, 展示以及回滚的时候再替换回来
func formatData(row *sql.Rows, columns *[]sql_util.Column) string {
	defer row.Close()
	values, err := scanner.ScanMap(row)
	if err != nil {
		slog.Errorf("format data scanMap rows failed : %s", err.Error())
		return ""
	}

	var resp string
	for _, rowMap := range values {
		rowStr := ""
		for _, column := range *columns {
			if fieldValue, ok := rowMap[column.Field]; ok {
				rowStr += splitFieldAsterisk + convertField(uint8ToString(fieldValue))
			} else {
				slog.Errorf("backup data format data error, column : %s not found . data : %v", column.Field, rowMap)
			}
		}
		if len(rowStr) >= 1 {
			rowStr = rowStr[1:]
		}
		resp += splitRowNumberSign + rowStr
	}
	if len(resp) >= 1 {
		resp = resp[1:]
	}
	return resp
}

func convertField(fieldStr string) string {
	if fieldStr == NUL {
		return replaceNUL
	}
	fieldStr = strings.ReplaceAll(fieldStr, splitRowNumberSign, replaceNumberSign)
	return strings.ReplaceAll(fieldStr, splitFieldAsterisk, replaceAsterisk)
}

func reverseConvertFields(str []string) (resp []string) {
	for _, v := range str {
		resp = append(resp, reverseConvertField(v))
	}
	return
}

func reverseConvertField(str string) string {
	if str == replaceNUL {
		return NUL
	}
	str = strings.ReplaceAll(str, replaceAsterisk, splitFieldAsterisk)
	return strings.ReplaceAll(str, replaceNumberSign, splitRowNumberSign)
}

func uint8ToString(inter interface{}) string {
	uint8Array, ok := inter.([]uint8)
	if !ok {
		slog.Errorf("uint8 to string error : received interface isn't a uint8 slice, data: %v", inter)
		return ""
	}

	var byteArray []byte
	for _, v := range uint8Array {
		byteArray = append(byteArray, byte(v))
	}
	return string(byteArray)
}

func isNum(str string) bool {
	for _, v := range str {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
