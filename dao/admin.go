package dao

import (
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
)

type AdminDaoImpl struct {
}

var Admin AdminDaoImpl

func (AdminDaoImpl) AddAdmin(admin *admin.DbInjectionAdmin) (int64, error) {
	err := GetDB().Create(admin).Error
	return admin.ID, err
}

func (AdminDaoImpl) GetAdmin(username string) (*admin.DbInjectionAdmin, error) {
	var admin admin.DbInjectionAdmin
	if err := GetDB().First(&admin, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}
