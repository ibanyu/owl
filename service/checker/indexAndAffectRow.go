package checker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/shawnfeng/sutil/slog"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/sql_util"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xsql/scanner"
)

type opType int

const (
	Update opType = iota + 1
	Delete
)

func IndexMach(info *task.DBInfo, sql string) (bool, error) {
	sql = strings.ToLower(sql)
	if !strings.Contains(sql, "where") {
		slog.Infof("sql check： sql not contain 'where' ")
		return false, errors.New("sql not contain 'where' ")
	}

	sql = strings.TrimSpace(sql)
	sqlAfterWhere := sql_util.GetSqlAfterWhere(sql)
	if operateDisableIndex(sqlAfterWhere) {
		return false, nil
	}
	keysInfo, err := getIndexInfo(info, sql)
	if err != nil {
		slog.Warnf("sql check, get index info err: %s", err.Error())
		return false, err
	}
	condition := sql_util.GetCondition(sqlAfterWhere)

	return indexMatchConditionOrdinal(keysInfo, condition) || indexMatchConditionAllValue(keysInfo, condition), nil
}

// 条件的值和某索引乱序匹配
func indexMatchConditionAllValue(keys *[]KeysInfo, condition []string) bool {
	keySplit := make(map[string][]KeysInfo)
	for _, v := range *keys {
		keySplit[v.KeyName] = append(keySplit[v.KeyName], v)
	}

	for _, oneKey := range keySplit {
		if allValEqual(oneKey, condition) {
			return true
		}
	}
	return false
}

func allValEqual(oneKey []KeysInfo, condition []string) bool {
	if len(oneKey) != len(condition) {
		return false
	}
	for _, keyCol := range oneKey {
		find := false
		for _, condition := range condition {
			if keyCol.ColumnName == condition {
				find = true
			}
		}
		if !find {
			return false
		}
	}
	return true
}

// 条件按顺序匹配索引
func indexMatchConditionOrdinal(keys *[]KeysInfo, condition []string) bool {
	for i, v := range condition {
		matchI := false
		for _, key := range *keys {
			if strings.ToLower(key.ColumnName) == v && key.SeqInIndex == i+1 {
				matchI = true
				break
			}
		}
		if !matchI {
			slog.Infof("sql condition not math, num : %d, condition: %s, index: %v", i, v, *keys)
			return false
		}
	}
	return true
}

type KeysInfo struct {
	KeyName    string `bdb:"Key_name"`
	SeqInIndex int    `bdb:"Seq_in_index"` // index,索引 即key，
	ColumnName string `bdb:"Column_name"`
}

// 获取索引信息
func getIndexInfo(info *task.DBInfo, sql string) (*[]KeysInfo, error) {
	tableName, err := sql_util.GetTableName(sql)
	if err != nil {
		return nil, err
	}
	indexSql := fmt.Sprintf("show index from %s", tableName)
	slog.Infof("build get index sql : %s", indexSql)

	res, err := info.DB.QueryContext(context.TODO(), indexSql)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	keysInfo := &[]KeysInfo{}
	err = scanner.Scan(res, keysInfo)
	if err != nil {
		return nil, err
	}
	slog.Infof("sql index info: %v, index sql: %s", *keysInfo, indexSql)
	return keysInfo, nil
}

//仅做导致索引失效的包含条件判断用
const (
	indexDisableConditionIn      = " in "
	indexDisableConditionNotIn   = " not in "
	indexDisableConditionLike    = " like "
	indexDisableConditionBetween = " between "
	indexDisableConditionOr      = " or "
	indexDisableConditionNull    = "null"
)

//检查会导致索引失效的操作,
// or， ！=、not in、<> ,  条件包含null， like 百分号开头， 表达式操作：+ - * / , 等号左边函数操作：substring、dateadd、year、 等（最后这个不考虑）
// 范围条件后面的索引也会失效， > < between 不在最后一列的情况
func operateDisableIndex(sqlAfterWhere string) bool {
	sqlAfterWhere = strings.ToLower(sqlAfterWhere)
	if strings.Contains(sqlAfterWhere, indexDisableConditionOr) ||
		strings.Contains(sqlAfterWhere, "!=") ||
		strings.Contains(sqlAfterWhere, indexDisableConditionNotIn) ||
		strings.Contains(sqlAfterWhere, "<>") ||
		strings.Contains(sqlAfterWhere, indexDisableConditionNull) {
		slog.Infof("sql contain 'or' or '!=' or '<>' or 'null' ")
		return true
	}

	if strings.Contains(sqlAfterWhere, "+") ||
		strings.Contains(sqlAfterWhere, "-") ||
		strings.Contains(sqlAfterWhere, "*") ||
		strings.Contains(sqlAfterWhere, "/") {
		slog.Infof(`sql contain '+' or '-' or '*' or '/' `)
		return true
	}

	//包含like，且%开头的检查，判断方式： 去掉空格，like 关键字往后查两个看是不是%
	sql := strings.ReplaceAll(sqlAfterWhere, " ", "")
	for {
		if index := strings.Index(sql, "like"); index > 0 {
			if len(sql) < index+5 {
				slog.Infof("sql check error : like have no value")
				//正确性检测与此无关,其他地方会做
				return true
			}
			if sql[index+5] == '%' {
				slog.Infof("sql contain 'like' , and begin as percent sign ")
				return true
			}
			sql = strings.Replace(sql, "like", "", 1)
			continue
		}
		break
	}

	resp := scopeOperateDisableIndex(sqlAfterWhere)
	if resp {
		slog.Infof("sql contains scope condition, disable index ")
	}
	return resp
}

//包含>,< ,>=, <= ,like,in, between ,且不是最后一列的检查。 是不是最后一列的判断，其后没有and；between 是有没有超过一个and
// 参数为已经转小写的sql
func scopeOperateDisableIndex(sqlAfterWhere string) bool {
	//同样的范围条件判断
	if strings.Count(sqlAfterWhere, ">") > 1 ||
		strings.Count(sqlAfterWhere, "<") > 1 ||
		strings.Count(sqlAfterWhere, ">=") > 1 ||
		strings.Count(sqlAfterWhere, "<=") > 1 ||
		strings.Count(sqlAfterWhere, indexDisableConditionLike) > 1 ||
		strings.Count(sqlAfterWhere, indexDisableConditionIn) > 1 ||
		strings.Count(sqlAfterWhere, indexDisableConditionBetween) > 1 {
		return true
	}

	switch {
	case strings.Contains(sqlAfterWhere, ">"):
		if strings.Contains(strings.Split(sqlAfterWhere, ">")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, "<"):
		if strings.Contains(strings.Split(sqlAfterWhere, "<")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, ">="):
		if strings.Contains(strings.Split(sqlAfterWhere, ">=")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, "<="):
		if strings.Contains(strings.Split(sqlAfterWhere, "<=")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, indexDisableConditionLike):
		if strings.Contains(strings.Split(sqlAfterWhere, "like")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, indexDisableConditionIn):
		if strings.Contains(strings.Split(sqlAfterWhere, "in")[1], "and") {
			return true
		}
	case strings.Contains(sqlAfterWhere, indexDisableConditionBetween):
		if strings.Count(strings.Split(sqlAfterWhere, "between")[1], "and") > 1 {
			return true
		}
	}

	return false
}

func getAffectRow(info *task.DBInfo, sql string) (int, error) {
	sql, err := dmlSqlToCount(sql)
	if err != nil {
		slog.Warnf("get dml convert select sql error : %s", err.Error())
		return -1, err
	}
	var count int64
	err = info.DB.QueryRow(sql).Scan(&count)
	if err != nil {
		slog.Warnf("get dml sql affect row error : %s", err.Error())
		return -1, err
	}

	return int(count), nil
}

func dmlSqlToCount(sql string) (string, error) {
	sqlLower := strings.ToLower(sql)
	if !strings.Contains(sqlLower, "where") {
		return "", errors.New("sql translate to count, not contain 'where'")
	}

	tableName, err := sql_util.GetTableName(sql)
	if err != nil {
		return "", err
	}
	countSql := fmt.Sprintf("select count(*) from %s where %s", tableName, sql_util.GetSqlAfterWhere(sql))
	slog.Infof("build count sql : %s", countSql)

	return countSql, nil
}
