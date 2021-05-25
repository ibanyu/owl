package auth

import (
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
)

type AdminAuthToolImpl struct {
}

var AdminAuthService ConfAuthToolImpl

func (AdminAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	admins, _, err := admin.ListAdmin(&service.Pagination{})
	if err != nil {
		return "", err
	}

	var resp string
	for i, v := range admins {
		if i == 0 {
			resp += v.Username
		} else {
			resp += "," + v.Username
		}
	}
	return resp, nil
}

func (AdminAuthToolImpl) IsDba(userName string) (isDba bool, err error) {
	admins, _, err := admin.ListAdmin(&service.Pagination{})
	if err != nil {
		return false, err
	}

	for _, v := range admins {
		if v.Username == userName {
			return true, nil
		}
	}

	return false, nil
}
