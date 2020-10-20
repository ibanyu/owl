package db_info

import (
	"database/sql"
	"fmt"
)

type Database struct {
	Database []string `json:"database"`
}

func GetDBConn(clusterName, dbName string) (*sql.DB, error) {
	cluster, err := clusterDao.GetClusterByName(clusterName)
	if err != nil{
		return nil, err
	}

	return sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, "test"))
}

// return dbs and mapping cluster
func ListAllDB() (map[string][]string, error) {
	clusters, err := ListCluster()
	if err != nil {
		return nil, err
	}

	resp := make(map[string][]string)
	for _, cluster := range clusters {
		dbs, err := ListDbByCluster(&cluster)
		if err != nil {
			return nil, err
		}

		for _, db := range dbs {
			resp[db] = append(resp[db], cluster.Name)
		}
	}

	return resp, nil
}

func ListDbByCluster(cluster *DbInjectionCluster) ([]string, error) {
	conn, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, "test"))
	if err != nil {
		return nil, fmt.Errorf("open db_info conn err: %s", err.Error())
	}

	rows, err := conn.Query("show databases;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs Database
	return dbs.Database, rows.Scan(&dbs)
}
