package wxpush

import "fmt"

type Error struct {
	Code int
	Msg  string
	Err  error
}

func (err Error) Error() string {
	return fmt.Sprintf("%d %s %s", err.Code, err.Msg, err.Err)
}
func NewError(code int, err error) error {
	var msg string
	switch code {
	case 1001:
		msg = "业务异常错误"
	case 1002:
		msg = "未认证"
	case 1003:
		msg = "签名错误"
	case 1004:
		msg = "接口不存在"
	case 1005:
		msg = "服务器内部错误"
	case 1006:
		msg = "和微信交互的过程中发生异常"
	case 1007:
		msg = "网络异常"
	case 1008:
		msg = "数据异常"
	case 1009:
		msg = "未知异常"
	default:
		msg = "未知异常"
	}
	return Error{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}
func NewBusinessError(err error) error {
	return NewError(1001, err)
}
