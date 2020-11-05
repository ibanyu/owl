package test

import (
	"github.com/gin-gonic/gin"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
	"net/http"
	"net/http/httptest"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/router"
)

var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZpc2giLCJwYXNzd29yZCI6ImZseSJ9.8Lur8qWsME-nI_TdFS7aGqAUa4sbup8Qf2Lb5Oikx1g"

const (
	userName = "fish"
	pwd = "fly"
)

func init()  {
	logger.InitLog()
	config.InitConfig("../../config/config.yml")
}

func serverRouter(req *http.Request) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := router.Router()
	respWriter := httptest.NewRecorder()
	router.ServeHTTP(respWriter, req)
	return respWriter
}