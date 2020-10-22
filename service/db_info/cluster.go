package db_info

type DbInjectionCluster struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	DefaultDB   string `json:"default_db" gorm:"column:default_db"`
	Addr        string `json:"addr" gorm:"column:addr"` //ip : port
	User        string `json:"user" gorm:"column:user"`
	Pwd         string `json:"pwd" gorm:"column:pwd"`

	Ct       int64  `json:"ct" gorm:"column:ct"`
	Ut       int64  `json:"ut" gorm:"column:ut"`
	Operator string `json:"operator" gorm:"column:operator"`
}

type ClusterDao interface {
	AddCluster(cluster *DbInjectionCluster) (int64, error)
	UpdateCluster(cluster *DbInjectionCluster) error
	DelCluster(id int64) error
	GetClusterByName(clusterName string) (*DbInjectionCluster, error)
	ListCluster() ([]DbInjectionCluster, error)
}

var clusterDao ClusterDao

func SetClusterDao(impl ClusterDao) {
	clusterDao = impl
}

func AddCluster(cluster *DbInjectionCluster) (int64, error) {
	return clusterDao.AddCluster(cluster)
}

func UpdateCluster(cluster *DbInjectionCluster) error {
	return clusterDao.UpdateCluster(cluster)
}

func DelCluster(id int64) error {
	return clusterDao.DelCluster(id)
}

func GetClusterByName(name string) (*DbInjectionCluster, error) {
	return clusterDao.GetClusterByName(name)
}

func ListCluster() ([]DbInjectionCluster, error) {
	return clusterDao.ListCluster()
}
