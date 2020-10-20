package code

import (
	"fmt"
)

const (
	Success = 0

	ParamInvalid   = 400
	UserNotFound   = 401
	PageNotFound   = 404
	InvalidOperate = 405

	InternalErr    = 500
	InternalNetErr = 501
)

type Err struct {
	ErrCode int    `json:"code"`
	Msg     string `json:"msg"`
}

func (e Err) Error() string {
	return e.Msg
}

func (e Err) Code() int {
	return e.ErrCode
}

func NewError(code int, msg string, args ...interface{}) Err {
	return Err{
		ErrCode: code,
		Msg:     fmt.Sprintf(msg, args...),
	}
}
