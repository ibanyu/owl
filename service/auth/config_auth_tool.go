package auth

import (
	"errors"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
)

type ConfAuthToolImpl struct {
}

var ConfAuthService ConfAuthToolImpl

func (ConfAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	for _, v := range config.Conf.Role.Conf.ReviewerRelation {
		for _, member := range v.Members {
			if userName == member {
				return v.Reviewer, nil
			}
		}
	}

	return "", errors.New("get reviewer from config failed, not found")
}

func (ConfAuthToolImpl) IsDba(userName string) (isDba bool, err error) {
	if len(config.Conf.Role.Conf.DBA) < 1 {
		return false, errors.New("dba members not config")
	}
	for _, v := range config.Conf.Role.Conf.DBA {
		if v == userName {
			return true, nil
		}
	}

	return false, nil
}
