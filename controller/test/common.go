package test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibanyu/owl/config"
	"github.com/ibanyu/owl/router"
	"github.com/ibanyu/owl/service"
	"github.com/ibanyu/owl/util/logger"
)

var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZpc2giLCJwYXNzd29yZCI6ImZseSJ9.8Lur8qWsME-nI_TdFS7aGqAUa4sbup8Qf2Lb5Oikx1g"

const (
	userName = "fish"
	pwd      = "fly"
)

func init() {
	logger.InitLog(".", "test.log", "debug")
	config.InitConfig("../../config/config.yml")
	service.SetClock(MockClock{})
}

func serverRouter(req *http.Request) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := router.Router()
	respWriter := httptest.NewRecorder()
	router.ServeHTTP(respWriter, req)
	return respWriter
}

type MockClock struct{}

func (MockClock) Now() time.Time { return time.Time{} }
func (MockClock) NowUnix() int64 { return 0 }
