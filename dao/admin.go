package dao

import (
	"github.com/ibanyu/owl/service"
	"github.com/ibanyu/owl/service/admin"
)

type AdminDaoImpl struct {
}

var Admin AdminDaoImpl

func (AdminDaoImpl) AddAdmin(admin *admin.OwlAdmin) (int64, error) {
	err := GetDB().Create(admin).Error
	return admin.ID, err
}

func (AdminDaoImpl) GetAdmin(username string) (*admin.OwlAdmin, error) {
	var admin admin.OwlAdmin
	if err := GetDB().First(&admin, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

func (AdminDaoImpl) ListAdmin(page *service.Pagination) ([]admin.OwlAdmin, int, error) {
	condition := "username like ?"

	page.Key = "%" + page.Key + "%"
	var count int
	if err := GetDB().Model(&admin.OwlAdmin{}).Where(condition,
		page.Key).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var admins []admin.OwlAdmin
	if err := GetDB().Order("ct desc").Offset(page.Offset).Limit(page.Limit).
		Find(&admins, condition, page.Key).Error; err != nil {
		return nil, 0, err
	}

	return admins, count, nil
}

func (AdminDaoImpl) DelAdmin(id int64) error {
	return GetDB().Where("id = ?", id).Delete(&admin.OwlAdmin{}).Error
}
