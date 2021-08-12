package dao

import (
	"github.com/ibanyu/owl/service/checker"
	"github.com/jinzhu/gorm"
)

type RuleDaoImpl struct {
}

var Rule RuleDaoImpl

func (RuleDaoImpl) ListAllStatus() ([]checker.DbInjectionRuleStatus, error) {
	var ruleStatus []checker.DbInjectionRuleStatus
	return ruleStatus, GetDB().Find(&ruleStatus).Error
}

func (RuleDaoImpl) UpdateRuleStatus(ruleStatus *checker.DbInjectionRuleStatus) error {
	err := GetDB().Where("name = ?", ruleStatus.Name).First(&checker.DbInjectionRuleStatus{}).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return GetDB().Create(ruleStatus).Error
	}
	if err != nil {
		return err
	}

	return GetDB().Model(ruleStatus).Where("name = ?", ruleStatus.Name).Update(ruleStatus).Error
}
