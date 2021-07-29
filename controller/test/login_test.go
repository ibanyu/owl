package test

import (
	"bytes"
	"encoding/json"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth"
	"gotest.tools/assert"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller/test/mock"
)

func injectAuthMock(t *testing.T) *mock.MockLoginChecker {
	ctl := gomock.NewController(t)

	mockAuth := mock.NewMockLoginChecker(ctl)
	auth.SetLoginService(mockAuth)
	return mockAuth
}

func TestLogin(t *testing.T) {
	mockAuth := injectAuthMock(t)

	mockAuth.EXPECT().Login(userName, pwd).Return(nil)
	//expect

	page := auth.Claims{
		Username: userName,
		Password: pwd,
	}
	pageByte, _ := json.Marshal(page)
	req, _ := http.NewRequest("POST", "/db-injection/login", bytes.NewBuffer(pageByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}
