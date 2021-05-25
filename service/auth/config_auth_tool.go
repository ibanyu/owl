package auth

import (
	"errors"
	"strings"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
)

type ConfAuthToolImpl struct {
}

var ConfAuthService ConfAuthToolImpl

func (ConfAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	return  strings.Join(config.Conf.Role.Conf.DBA, ","), nil
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
