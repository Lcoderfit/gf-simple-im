package model

import "github.com/gogf/gf/net/ghttp"

const (
	// 存储上下文变量的键名
	ContextKey = "ContextKey"
)

// 自定义上下文
type Context struct {
	Session *ghttp.Session // 当前Session管理对象
	User    *ContextUser   // 上下文用户信息
}

// 自定义用户上下文
type ContextUser struct {
	Id       uint   // 用户ID
	Passport string // 用户账号
	Nickname string // 用户昵称
}
