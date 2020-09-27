package advisor

import (
	"context"
	"fmt"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/mysql"
	"github.com/shawnfeng/sutil/slog"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xsql/scanner"
)

const keyJoinChar = "+"

func isSubKey(keyA, keyB string) bool {
	var short, long string
	if len(keyA) < len(keyB) {
		short, long = keyA, keyB
	} else {
		short, long = keyB, keyA
	}

	subKeys := strings.Split(long, keyJoinChar)
	for _, v := range subKeys {
		if v == short {
			return true
		}
	}
	return false
}

//取出数据的顺序同建表语句的顺序
type Column struct {
	Field   string `bdb:"Field"`
	Type    string `bdb:"Type"`
	Null    string `bdb:"Null"`
	Key     string `bdb:"Key"`
	Default string `bdb:"Default"`
	Extra   string `bdb:"Extra"`
}

func GetTableColumn(tableName string, info *task.DBInfo) (*[]Column, error) {
	sql := fmt.Sprintf("show columns from %s ", tableName)
	db := info.DB

	rows, err := db.QueryContext(context.TODO(), sql)
	if err != nil {
		slog.Warnf("get table column info err:%v, sql content: %s", err, sql)
		return nil, err
	}
	var column []Column
	err = scanner.Scan(rows, &column)
	slog.Infof("table column is: %v", column)
	return &column, err
}

type WriterBuffer struct {
	Condition strings.Builder
}

func (v *WriterBuffer) Write(p []byte) (n int, err error) {
	v.Condition.Write(p)
	return
}

var Parser *parser.Parser = parser.New()

//返回值为where后面的部分，不包括where
func GetSqlAfterWhere(sql string) string {
	stmtNodes, _, err := Parser.Parse(sql, "", "")
	if err != nil {
		slog.Errorf("get sql condition err: %s", err.Error())
	}

	writer := &WriterBuffer{}
	ctx := format.NewRestoreCtx(format.RestoreStringSingleQuotes|format.RestoreKeyWordLowercase, writer)
	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.UpdateStmt:
			node.Where.Restore(ctx)
		case *ast.DeleteStmt:
			node.Where.Restore(ctx)
		default:
			slog.Errorf("get sql after err, not supported type. sql: %s", sql)
		}
	}
	return writer.Condition.String()
}

// add `` to key words
// todo 有瑕疵，同时有关键字和value内包含关键字内容的情况（eg: xxxx  where `index` = 'index hello'; ），会出问题。待更换更好的方式。
func HandelKeyWorldForCondition(originWhere string) string {
	targetWhere := originWhere
	columns := getCondition(originWhere)
	for _, v := range columns {
		if isKeyWord(v) {
			targetWhere = strings.ReplaceAll(targetWhere, v, fmt.Sprintf("`%s`", v))
		}
	}

	return targetWhere
}

// tidb 暂时不支持'SELECT * FROM mysql.help_keyword;', 支持后可以替换一下
var keyWords = map[string]struct{}{
	"ADD": {}, "ADMIN": {}, "ALL": {}, "ALTER": {}, "ANALYZE": {}, "AND": {}, "AS": {}, "ASC": {},

	"BETWEEN": {}, "BIGINT": {}, "BINARY": {}, "BLOB": {}, "BOTH": {}, "BUCKETS": {}, "BUILTINS": {}, "BY": {},

	"CANCEL": {}, "CASCADE": {}, "CASE": {}, "CHANGE": {}, "CHAR": {}, "CHARACTER": {}, "CHECK": {}, "CMSKETCH": {}, "COLLATE": {},
	"COLUMN": {}, "CONSTRAINT": {}, "CONVERT": {}, "CREATE": {}, "CROSS": {}, "CURRENT_DATE": {}, "CURRENT_ROLE": {}, "CURRENT_TIME": {}, "CURRENT_TIMESTAMP": {},
	"CURRENT_USER": {},

	"DATABASE": {}, "DATABASES": {}, "DAY_HOUR": {}, "DAY_MICROSECOND": {}, "DAY_MINUTE": {}, "DAY_SECOND": {}, "DDL": {}, "DECIMAL": {}, "DEFAULT": {},
	"DELAYED": {}, "DELETE": {}, "DEPTH": {}, "DESC": {}, "DESCRIBE": {}, "DISTINCT": {}, "DISTINCTROW": {}, "DIV": {}, "DOUBLE": {},
	"DRAINER": {}, "DROP": {}, "DUAL": {},

	"ELSE": {}, "ENCLOSED": {}, "ESCAPED": {}, "EXCEPT": {}, "EXISTS": {}, "EXPLAIN": {},

	"FALSE": {}, "FLOAT": {}, "FOR": {}, "FORCE": {}, "FOREIGN": {}, "FROM": {}, "FULLTEXT": {},

	"GENERATED": {}, "GRANT": {}, "GROUP": {}, "HAVING": {}, "HIGH_PRIORITY": {}, "HOUR_MICROSECOND": {}, "HOUR_MINUTE": {}, "HOUR_SECOND": {}, "IF": {},

	"IGNORE": {}, "IN": {}, "INDEX": {}, "INFILE": {}, "INNER": {}, "INSERT": {}, "INT": {}, "INT1": {}, "INT2": {},
	"INT3": {}, "INT4": {}, "INT8": {}, "INTEGER": {}, "INTERVAL": {}, "INTO": {}, "IS": {},

	"JOB": {}, "JOBS": {}, "JOIN": {}, "KEY": {}, "KEYS": {}, "KILL": {},

	"LEADING": {}, "LEFT": {}, "LIKE": {}, "LIMIT": {}, "LINEAR": {}, "LINES": {}, "LOAD": {}, "LOCALTIME": {}, "LOCALTIMESTAMP": {},
	"LOCK": {}, "LONG": {}, "LONGBLOB": {}, "LONGTEXT": {}, "LOW_PRIORITY": {},

	"MATCH": {}, "MAXVALUE": {}, "MEDIUMBLOB": {}, "MEDIUMINT": {}, "MEDIUMTEXT": {}, "MINUTE_MICROSECOND": {}, "MINUTE_SECOND": {}, "MOD": {},

	"NATURAL": {}, "NODE_ID": {}, "NODE_STATE": {}, "NOT": {}, "NO_WRITE_TO_BINLOG": {}, "NULL": {}, "NUMERIC": {},

	"ON": {}, "OPTIMISTIC": {}, "OPTIMIZE": {}, "OPTION": {}, "OPTIONALLY": {}, "OR": {}, "ORDER": {}, "OUTER": {}, "OUTFILE": {},

	"PARTITION": {}, "PESSIMISTIC": {}, "PRECISION": {}, "PRIMARY": {}, "PROCEDURE": {}, "PUMP": {},

	"RANGE": {}, "READ": {}, "REAL": {}, "REFERENCES": {}, "REGEXP": {}, "REGION": {}, "REGIONS": {}, "RELEASE": {}, "RENAME": {},
	"REPEAT": {}, "REPLACE": {}, "REQUIRE": {}, "RESTRICT": {}, "REVOKE": {}, "RIGHT": {}, "RLIKE": {}, "ROW": {},

	"SAMPLES": {}, "SECOND_MICROSECOND": {}, "SELECT": {}, "SET": {}, "SHOW": {}, "SMALLINT": {}, "SPATIAL": {}, "SPLIT": {}, "SQL": {},
	"SQL_BIG_RESULT": {}, "SQL_CALC_FOUND_ROWS": {}, "SQL_SMALL_RESULT": {}, "SSL": {}, "STARTING": {}, "STATS": {}, "STATS_BUCKETS": {}, "STATS_HEALTHY": {}, "STATS_HISTOGRAMS": {},
	"STATS_META": {}, "STORED": {}, "STRAIGHT_JOIN": {},

	"TABLE": {}, "TERMINATED": {}, "THEN": {}, "TIDB": {}, "TIFLASH": {}, "TINYBLOB": {}, "TINYINT": {}, "TINYTEXT": {}, "TO": {},
	"TOPN": {}, "TRAILING": {}, "TRIGGER": {}, "TRUE": {},

	"UNION": {}, "UNIQUE": {}, "UNLOCK": {}, "UNSIGNED": {}, "UPDATE": {}, "USAGE": {}, "USE": {}, "USING": {}, "UTC_DATE": {},
	"UTC_TIME": {}, "UTC_TIMESTAMP": {},

	"VALUES": {}, "VARBINARY": {}, "VARCHAR": {}, "VARCHARACTER": {}, "VARYING": {}, "VIRTUAL": {},

	"WHEN": {}, "WHERE": {}, "WIDTH": {}, "WITH": {}, "WRITE": {}, "XOR": {}, "YEAR_MONTH": {}, "ZEROFILL": {},
}

func isKeyWord(name string) bool {
	name = strings.ToUpper(name)
	if _, ok := keyWords[name]; ok {
		return true
	}
	return false
}

const varcharLimit = 1024

func varcharLengthTooLong(nodes []ast.StmtNode) bool {
	for _, tiStmt := range nodes {
		switch node := tiStmt.(type) {
		case *ast.AlterTableStmt:
			for _, v := range node.Specs {
				for _, col := range v.NewColumns {
					if col.Tp.Tp == mysql.TypeVarchar && col.Tp.Flen > varcharLimit {
						return true
					}
				}
			}
		case *ast.CreateTableStmt:
			for _, v := range node.Cols {
				if v.Tp.Tp == mysql.TypeVarchar && v.Tp.Flen > varcharLimit {
					return true
				}
			}
		default:
			return false
		}
	}
	return false
}

func singlePrimaryKeyIsInt(nodes []ast.StmtNode) bool {
	for _, tiStmt := range nodes {
		switch node := tiStmt.(type) {
		case *ast.CreateTableStmt:
			primaryKeyColName := ""
			for _, v := range node.Constraints {
				if v.Tp == ast.ConstraintPrimaryKey && len(v.Keys) == 1 {
					primaryKeyColName = v.Keys[0].Column.Name.L
				}
			}

			if primaryKeyColName == "" {
				return false
			}

			for _, v := range node.Cols {
				if v.Name.Name.L == primaryKeyColName &&
					(v.Tp.Tp == mysql.TypeShort ||
						v.Tp.Tp == mysql.TypeLong ||
						v.Tp.Tp == mysql.TypeTiny ||
						v.Tp.Tp == mysql.TypeInt24) {
					return true
				}
			}
		default:
			return false
		}
	}
	return false
}
