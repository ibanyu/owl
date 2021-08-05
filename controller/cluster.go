package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"
)

func ListDB(ctx *gin.Context) Resp {
	f := "ListDB()-->"

	dbInfo, err := db_info.ListAllDB()
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s,list db failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{Data: dbInfo}
}

func ListCluster(ctx *gin.Context) Resp {
	f := "ListCluster()-->"

	cluster, err := db_info.ListClusterForUI()
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s,list cluster failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{Data: ListData{
		Items:  cluster,
		Total:  len(cluster),
		More:   false,
		Offset: 0,
	}}
}

func AddCluster(ctx *gin.Context) Resp {
	f := "AddCluster()-->"

	var cluster db_info.DbInjectionCluster
	if err := ctx.BindJSON(&cluster); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	id, err := db_info.AddCluster(&cluster)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s,add cluster failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{Data: id}
}

func UpdateCluster(ctx *gin.Context) Resp {
	f := "UpdateCluster()-->"

	var cluster db_info.DbInjectionCluster
	if err := ctx.BindJSON(&cluster); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	if err := db_info.UpdateCluster(&cluster); err != nil {
		return Resp{Message: fmt.Sprintf("%s,update cluster failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{}
}

func DelCluster(ctx *gin.Context) Resp {
	f := "DelCluster()-->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Resp{
			Message: fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr),
			Code:    code.ParamInvalid,
		}
	}

	if err := db_info.DelCluster(id); err != nil {
		return Resp{Message: fmt.Sprintf("%s,del cluster failed :%s ", f, err.Error()), Code: code.ParamInvalid}
	}

	return Resp{}
}
