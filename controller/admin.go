package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
)

func AddAdmin(ctx *gin.Context) Resp {
	f := "AddAdmin()-->"
	var adminParam admin.DbInjectionAdmin
	if err := ctx.BindJSON(&adminParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}

	adminParam.Creator = ctx.MustGet("user").(string)
	id, err := admin.AddAdmin(&adminParam)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s, add admin failed :%s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: struct {
		Id int64 `json:"id"`
	}{id}}
}
