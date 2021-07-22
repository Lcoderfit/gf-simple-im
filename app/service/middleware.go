package service

import (
	"gf-simple-im/app/model"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

// 中间件管理服务
var Middleware = middlewareService{}

type middlewareService struct{}

// 自定义上下文对象
func (s *middlewareService) Ctx(r *ghttp.Request) {
	// 每个用户都有不同的session和user信息，所以每次都需要重新创建一个
	customCtx := &model.Context{
		Session: r.Session,
	}
	// 将自定义context作为系统的context
	Context.Init(r, customCtx)
	// 当尚未登录时，user应该是为nil的，所以会直接执行后面的中间件
	// 登录之后，则会再登录接口中设置customCtx.User和customCtx.Session
	// 此时再去访问其他接口，再次经过这个中间件时，均会带上session和user信息
	if user := Session.GetUser(r.Context()); user != nil {
		customCtx.User = &model.ContextUser{
			Id:       user.Id,
			Passport: user.Passport,
			Nickname: user.Nickname,
		}
	}
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// 鉴权中间件，只有登录成功之后才能通过
func (s *middlewareService) Auth(r *ghttp.Request) {
	if User.IsSignedIn(r.Context()) {
		r.Middleware.Next()
	} else {
		// 如果没有权限，则禁止访问
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

// 允许接口跨域请求
func (s *middlewareService) CORS(r *ghttp.Request) {
	// 设置默认的跨域选项，即允许任意的跨域请求
	r.Response.CORSDefault()
	// 调用下一个工作流处理程序
	r.Middleware.Next()
}
