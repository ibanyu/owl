package router

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

const (
	test    = 0
	product = 1
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/db-injection/login", HandlerWrapper(controller.Login))

	r.Use(controller.AuthorizeJWT())

	r.POST("/db-injection/user/role", HandlerWrapper(controller.RoleGet))

	task := r.Group("/db-injection/task")
	{
		task.POST("/add", HandlerWrapper(controller.AddTask))
		task.POST("/reviewer/list", HandlerWrapper(controller.ListReviewerTask))
		task.POST("/exec/list", HandlerWrapper(controller.ListExecTask))
		task.POST("/history", HandlerWrapper(controller.ListHistoryTask))
		task.POST("/get", HandlerWrapper(controller.GetTask))
		task.POST("/update", HandlerWrapper(controller.UpdateTask))
	}

	r.POST("/db-injection/rule/list", HandlerWrapper(controller.LisRule))
	r.POST("/db-injection/db/list", HandlerWrapper(controller.ListDB))

	backup := r.Group("/db-injection/backup")
	{
		backup.POST("/data", HandlerWrapper(controller.ListRollbackData))
		backup.POST("/rollback", HandlerWrapper(controller.Rollback))
	}

	r.Use(controller.OnlyDbaOrAdmin())

	r.POST("/db-injection/rule/update", HandlerWrapper(controller.UpdateRuleStatus))

	cluster := r.Group("/db-injection/cluster")
	{
		cluster.POST("/list", HandlerWrapper(controller.ListCluster))
		cluster.POST("/update", HandlerWrapper(controller.UpdateCluster))
		cluster.POST("/del", HandlerWrapper(controller.DelCluster))
		cluster.POST("/add", HandlerWrapper(controller.AddCluster))
	}

	admin := r.Group("/db-injection/admin")
	{
		admin.POST("/add", HandlerWrapper(controller.AddAdmin))
		admin.POST("/list", HandlerWrapper(controller.ListAdmin))
		admin.POST("/del", HandlerWrapper(controller.DelAdmin))
	}

	return r
}

func Run() {
	if config.Conf.Server.Env == product {
		gin.SetMode(gin.ReleaseMode)
	}

	r := Router()
	r.Use(gin.Recovery())

	endless.ListenAndServe(fmt.Sprintf(":%s", config.Conf.Server.Port), r)
}

type LogicHandler func(c *gin.Context) controller.Resp

func HandlerWrapper(handler LogicHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		writeResp(getHandleInfo(handler, c), c)
	}
}

func getHandleInfo(handler LogicHandler, c *gin.Context) controller.Resp {
	resp := handler(c)
	if resp.Message == "" && resp.Code == code.Success {
		resp.Message = "success"
	}

	if resp.Code != code.Success {
		logger.Errorf(resp.Message)
	}

	fName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	user, err := c.Cookie("user")
	if err != nil {
		logger.Warnf("%s get user in wrapper failed", fName)
	}

	if resp.Code == code.Success {
		logger.Infof("%s exec %s success", user, fName)
	} else {
		logger.Infof("%s exec %s failed", user, fName)
	}

	return resp
}

func writeResp(resp controller.Resp, c *gin.Context) {
	c.JSON(200, resp)
}
