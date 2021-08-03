package auth

import (
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
)

type AdminAuthToolImpl struct {
}

var AdminAuthService AdminAuthToolImpl

func (AdminAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	admins, _, err := admin.ListAdmin(&service.Pagination{Limit: 10})
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
	return admin.IsAdmin(userName)
}
