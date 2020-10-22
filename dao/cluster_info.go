package dao

import "gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"

type ClusterImpl struct {
}

var Cluster ClusterImpl

func (ClusterImpl) AddCluster(cluster *db_info.DbInjectionCluster) (int64, error) {
	err := GetDB().Create(cluster).Error
	return cluster.ID, err
}

func (ClusterImpl) UpdateCluster(cluster *db_info.DbInjectionCluster) error {
	return GetDB().Model(cluster).Where("id = ?", cluster.ID).Update(cluster).Error
}

func (ClusterImpl) DelCluster(id int64) error {
	return GetDB().Where("id = ?", id).Delete(&db_info.DbInjectionCluster{}).Error
}

func (ClusterImpl) GetClusterByName(name string) (*db_info.DbInjectionCluster, error) {
	var cluster db_info.DbInjectionCluster
	return &cluster, GetDB().First(&cluster, "name = ?", name).Error
}

func (ClusterImpl) ListCluster() ([]db_info.DbInjectionCluster, error) {
	var clusters []db_info.DbInjectionCluster
	return clusters, GetDB().Find(&clusters).Error
}
