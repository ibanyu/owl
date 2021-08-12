package test

import (
	"encoding/json"
	"gotest.tools/assert"
	"net/http"
	"testing"

	"github.com/ibanyu/owl/controller"
)

func TestListRule(t *testing.T) {
	req, _ := http.NewRequest("POST", "/db-injection/rule/list", nil)
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}
