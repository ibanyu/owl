package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service"
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

func ListAdmin(ctx *gin.Context) Resp {
	f := "ListAdmin() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	page.Operator = ctx.MustGet("user").(string)
	admins, count, err := admin.ListAdmin(&page)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s: list admin err: %s", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{Data: ListData{
		Items:  admins,
		Total:  count,
		More:   count > page.Offset+page.Limit,
		Offset: page.Offset,
	}}
}

func DelAdmin(ctx *gin.Context) Resp {
	f := "DelAdmin()-->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Resp{
			Message: fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr),
			Code:    code.ParamInvalid,
		}
	}

	if err := admin.DelAdmin(id); err != nil {
		return Resp{Message: fmt.Sprintf("%s,del admin failed :%s ", f, err.Error()), Code: code.InternalErr}
	}

	return Resp{}
}
