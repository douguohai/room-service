package entity

import (
	"fmt"
	"github.com/douguohai/room-service/base"
)

// Result 统一返回结果
type Result struct {
	Code int         `json:"code"` //状态码
	Msg  string      `json:"msg"`  //状态信息
	Data interface{} `json:"data"` //具体业务数据
}

// Success Result
func Success(msg string) Result {
	return Result{
		Msg:  msg,
		Code: base.SUCCESSCode,
	}
}

// Fail Result
func Fail(msg string) Result {
	return Result{
		Msg:  msg,
		Code: base.FailCode,
	}
}

// ToString Result
func (result *Result) ToString() string {
	return fmt.Sprintf("[%v,%v,%v]", result.Code, result.Msg, result.Data)
}
