package admin

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type DbInjectionAdmin struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Username    string `json:"username" gorm:"username"`
	Description string `json:"description" gorm:"column:description"`

	Ct      int64  `json:"ct" gorm:"column:ct"`
	Creator string `json:"creator" gorm:"creator"`
}

type AdminDao interface {
	AddAdmin(admin *DbInjectionAdmin) (int64, error)
	GetAdmin(username string) (*DbInjectionAdmin, error)
}

var adminDao AdminDao

func SetAdminDao(impl AdminDao) {
	adminDao = impl
}

func AddAdmin(admin *DbInjectionAdmin) (int64, error) {
	// add admin
	return adminDao.AddAdmin(admin)
}

func IsAdmin(username string) (bool, error) {
	_, err := adminDao.GetAdmin(username)
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get admin %s err", username)
	}

	return true, nil
}
