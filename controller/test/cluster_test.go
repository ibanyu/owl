package test

import (
	"bytes"
	"encoding/json"
	"gotest.tools/assert"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/controller/test/mock"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/db_info"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util"
)

func injectClusterMock(t *testing.T) *mock.MockClusterDao {
	ctl := gomock.NewController(t)
	injectAuthTool(t, ctl)

	mockClusterDao := mock.NewMockClusterDao(ctl)
	db_info.SetClusterDao(mockClusterDao)
	return mockClusterDao
}

func TestUpdateCluster(t *testing.T) {
	mockClusterDao := injectClusterMock(t)

	clusterIns := &db_info.DbInjectionCluster{
		ID:          777,
		Addr:        "db-cluster-addr",
		User:        "user",
		Pwd:         "db-cluster-pwd",
		Description: "this is a big db",
		Ct:          time.Now().Unix(),
	}

	wantCluster := *clusterIns
	pwdByte, err := util.AesCrypto([]byte(clusterIns.Pwd))
	assert.NilError(t, err)

	wantCluster.Pwd = util.StringifyByteDirectly(pwdByte)

	mockClusterDao.EXPECT().UpdateCluster(&wantCluster).Return(nil)

	clusterByte, _ := json.Marshal(clusterIns)
	req, _ := http.NewRequest("POST", "/db-injection/cluster/update", bytes.NewBuffer(clusterByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestDelCluster(t *testing.T) {
	mockClusterDao := injectClusterMock(t)

	mockClusterDao.EXPECT().DelCluster(int64(777)).Return(nil)

	req, _ := http.NewRequest("POST", "/db-injection/cluster/del?id=777", nil)
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestAddCluster(t *testing.T) {
	mockClusterDao := injectClusterMock(t)

	clusterIns := &db_info.DbInjectionCluster{
		Addr:        "db-cluster-addr",
		User:        "user",
		Pwd:         "db-cluster-pwd",
		Description: "this is a big db",
		Ct:          time.Now().Unix(),
	}

	wantCluster := *clusterIns
	pwdByte, err := util.AesCrypto([]byte(clusterIns.Pwd))
	assert.NilError(t, err)

	wantCluster.Pwd = util.StringifyByteDirectly(pwdByte)

	mockClusterDao.EXPECT().AddCluster(&wantCluster).Return(int64(1), nil)

	clusterByte, _ := json.Marshal(clusterIns)
	req, _ := http.NewRequest("POST", "/db-injection/cluster/add", bytes.NewBuffer(clusterByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestListCluster(t *testing.T) {
	mockClusterDao := injectClusterMock(t)

	ClusterIns := db_info.DbInjectionCluster{}
	mockClusterDao.EXPECT().ListCluster().Return([]db_info.DbInjectionCluster{ClusterIns}, nil)

	req, _ := http.NewRequest("POST", "/db-injection/cluster/list", nil)
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}
