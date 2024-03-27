package tools

import "fmt"

var (
	OK             = ECode{Code: 0, Message: "成功"}
	Success        = ECode{Code: 1, Message: "投票完成"}
	NotLogin       = ECode{Code: 10001, Message: "用户未登录"}
	ParamErr       = ECode{Code: 10002, Message: "参数错误"}
	UserErr        = ECode{Code: 10003, Message: "账号或密码错误"}
	UserExits      = ECode{Code: 10004, Message: "用户已经存在"}
	PassWordErr    = ECode{Code: 10005, Message: "两次密码不同"}
	PassWordNum    = ECode{Code: 10006, Message: "密码不能为纯数字"}
	UserCreateFail = ECode{Code: 10007, Message: "新用户创建失败！"}
	UserOver       = ECode{Code: 10008, Message: "用户名或密码大于8位小于16位"}
	VoteAlready    = ECode{Code: 10009, Message: "投票已存在"}
	VoteExits      = ECode{Code: 10010, Message: "您已经投过票"}
	CacheErr       = ECode{Code: 10011, Message: "验证码校验失败！"}
	DelFail        = ECode{Code: 10012, Message: "删除失败"}
	VoteNOAlready  = ECode{Code: 10013, Message: "该投票不存在"}
)

type ECode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("code:%d,message:%s", e.Code, e.Message)
}
