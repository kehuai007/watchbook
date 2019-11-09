package wxpush

import "errors"

// SendResult
type SendResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Uid       string `json:"uid"`
		TopicId   string `json:"topicId"`
		MessageId int    `json:"messageId"`
		Code      int    `json:"code"`
		Status    string `json:"status"`
	} `json:"data"`
	Success bool `json:"success"`
}

// Error 判断结果是否异常
func (result *SendResult) Error() error {
	if result.Success {
		return nil
	}
	if result.Code == 1000 {
		return nil
	}
	if result.Code == 0 {
		return NewBusinessError(errors.New("解析返回数据失败"))
	}
	return NewError(result.Code, errors.New(result.Msg))
}

type QueryMsgResult struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    int    `json:"data"`
	Success bool   `json:"success"`
}

func (result *QueryMsgResult) Error() error {
	if result.Success {
		return nil
	}
	if result.Code == 1000 {
		return nil
	}
	if result.Code == 0 {
		return NewBusinessError(errors.New("解析返回数据失败"))
	}
	return NewError(result.Code, errors.New(result.Msg))
}

type Users struct {
	Id         int    `json:"id"`
	Uid        string `json:"uid"`
	OpenId     string `json:"openId"`
	NickName   string `json:"nickName"`
	HeadImg    string `json:"headImg"`
	Sex        string `json:"sex"`
	ActiveTime string `json:"activeTime"`
	LastTime   string `json:"lastTime"`
	Enable     bool   `json:"enable"`
	Subscribe  string `json:"subscribe"`
	CreateTime int64  `json:"create_time"`
}

type QueryUsersResult struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		Total    int     `json:"total"`
		Page     int     `json:"page"`
		PageSize int     `json:"pageSize"`
		Records  []Users `json:"records"`
	} `json:"data"`
}

func (result QueryUsersResult) Error() error {
	if result.Success {
		return nil
	}
	if result.Code == 1000 {
		return nil
	}
	if result.Code == 0 {
		return NewBusinessError(errors.New("解析返回数据失败"))
	}
	return NewError(result.Code, errors.New(result.Msg))
}
