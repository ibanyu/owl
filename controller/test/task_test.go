package test

import (
	"bytes"
	"encoding/json"
	"gotest.tools/assert"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ibanyu/owl/controller"
	"github.com/ibanyu/owl/controller/test/mock"
	"github.com/ibanyu/owl/service"
	"github.com/ibanyu/owl/service/task"
)

func injectAuthTool(ctl *gomock.Controller) {
	mockAuth := mock.NewMockauthTools(ctl)
	mockAuth.EXPECT().IsDba(userName).Return(true, nil)
	mockAuth.EXPECT().GetReviewer(userName).Return(userName, nil)
	task.SetAuthTools(mockAuth)
}

func injectTaskMock(t *testing.T) *mock.MockTaskDao {
	ctl := gomock.NewController(t)
	injectAuthTool(ctl)

	mockTaskDao := mock.NewMockTaskDao(ctl)
	task.SetTaskDao(mockTaskDao)
	return mockTaskDao
}

func TestUpdateTask(t *testing.T) {
	mockTaskDao := injectTaskMock(t)

	taskIns := &task.DbInjectionTask{
		ID:       707,
		Status:   task.CheckPass,
		Executor: userName,
		Ut:       time.Now().Unix(),
		Action:   "progress",
	}

	taskWant := *taskIns
	taskWant.Action = ""
	taskWant.Status = task.ReviewPass

	mockTaskDao.EXPECT().UpdateTask(&taskWant).Return(nil)
	mockTaskDao.EXPECT().GetTask(taskIns.ID).Return(taskIns, nil)

	taskByte, _ := json.Marshal(taskIns)
	req, _ := http.NewRequest("POST", "/db-injection/task/update", bytes.NewBuffer(taskByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestGetTask(t *testing.T) {
	mockTaskDao := injectTaskMock(t)

	taskIns := &task.DbInjectionTask{
		Status:   task.CheckPass,
		Creator:  userName,
		Reviewer: userName,
		Ct:       time.Now().Unix(),
	}
	mockTaskDao.EXPECT().GetTask(int64(777)).Return(taskIns, nil)

	req, _ := http.NewRequest("POST", "/db-injection/task/get?id=777", nil)
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestAddTask(t *testing.T) {
	mockTaskDao := injectTaskMock(t)

	taskIns := &task.DbInjectionTask{
		Status:   task.CheckPass,
		Creator:  userName,
		Reviewer: userName,
		Ct:       time.Now().Unix(),
	}
	mockTaskDao.EXPECT().AddTask(taskIns).Return(int64(1), nil)

	pageByte, _ := json.Marshal(taskIns)
	req, _ := http.NewRequest("POST", "/db-injection/task/add", bytes.NewBuffer(pageByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestListTask(t *testing.T) {
	mockTaskDao := injectTaskMock(t)

	taskIns := task.DbInjectionTask{}
	page := service.Pagination{
		Offset:   5,
		Limit:    10,
		Operator: "fish",
	}
	mockTaskDao.EXPECT().ListTask(&page, true, task.ExecStatus()).Return([]task.DbInjectionTask{taskIns}, 1, nil)

	pageByte, _ := json.Marshal(page)
	req, _ := http.NewRequest("POST", "/db-injection/task/exec/list", bytes.NewBuffer(pageByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}

func TestListHistoryTask(t *testing.T) {
	mockTaskDao := injectTaskMock(t)

	taskIns := task.DbInjectionTask{}
	page := service.Pagination{
		Offset:   5,
		Limit:    10,
		Operator: "fish",
	}
	mockTaskDao.EXPECT().ListHistoryTask(&page, true).Return([]task.DbInjectionTask{taskIns}, 1, nil)

	pageByte, _ := json.Marshal(page)
	req, _ := http.NewRequest("POST", "/db-injection/task/history", bytes.NewBuffer(pageByte))
	req.Header.Set("token", token)
	respWriter := serverRouter(req)
	assert.Equal(t, 200, respWriter.Code)

	resp := &controller.Resp{}
	json.Unmarshal(respWriter.Body.Bytes(), resp)

	assert.Equal(t, 0, resp.Code)
}
