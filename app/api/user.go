package api

import (
	"gf-simple-im/app/model"
	"gf-simple-im/app/service"
	"gf-simple-im/library/response"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

// 用户API管理对象
var User = new(userApi)

type userApi struct{}

// true表示啥

// @summary 用户注册接口
// @tags 用户服务
// @produce json
// @param entity body model.UserApiSignUpReq true "注册请求"
// @success 200 {object} response.JsonResponse "执行结果"
// @router /user/signup [POST]
func (a *userApi) SignUp(r *ghttp.Request) {
	var (
		apiReq     *model.UserApiSignUpReq
		serviceReq *model.UserServiceSignUpReq
	)
	if err := r.ParseForm(&apiReq); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	// 将apiReq转换为serviceReq, 第二个参数传入指针，是因为需要修改serviceReq的值
	if err := gconv.Struct(apiReq, &serviceReq); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := service.User.SignUp(serviceReq); err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	response.JsonExit(r, 0, "ok")
}

// @summary 用户登录接口
// @tags 用户服务
// @produce json
// @param passport formData string true "用户账号"
// @param password formData string true "用户密码"
// @success 200 {object} response.JsonResponse "执行结果"
// @router /user/signin [POST]
func (a *userApi) SignIn(r *ghttp.Request) {
	var (
		data *model.UserApiSignInReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := service.User.SignIn(r.Context(), data.Passport, data.Password); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "ok")
}

// @summary 判断用户是否已经登录
// @tags 用户服务
// @produce json
// @success 200 {boject} response.JsonResponse "执行结果:`true/false`"
// @router /user/issignin [GET]
func (a *userApi) IsSignIn(r *ghttp.Request) {
	response.JsonExit(r, 0, "", service.User.IsSignedIn(r.Context()))
}

// @summary 用户注销/退出接口
// @tags 用户服务
// @produce json
// @success 200 {object} response.JsonResponse "执行结果, 1: 未登录"
// @router /user/signout [GET]
func (a *userApi) SignOut(r *ghttp.Request) {
	if err := service.User.SignOut(r.Context()); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "ok")
}

// @summary 检查用户账号接口(唯一性检测)
// @tags 用户服务
// @produce json
// @param passport query string true "用户帐号"
// @success 200 {object} response.JsonResponse "执行结果: `true/false`"
// @router /user/checkpassport [GET]
func (a *userApi) CheckPassport(r *ghttp.Request) {
	var (
		data *model.UserApiCheckPassportReq
	)
	// Parse将请求参数转换为golang结构体，并且会根据v标签进行字段验证
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data.Passport != "" && !service.User.CheckPassport(data.Passport) {
		response.JsonExit(r, 1, "账号已存在", false)
	}
	response.JsonExit(r, 0, "", true)
}

// @summary 检测用户昵称接口（唯一性校验）
// @tags 用户服务
// @produce json
// @param nickname query string true "用户昵称"
// @success 200 {object} response.JsonResponse "执行结果: `true/false`"
// @router /user/checknickname [GET]
func (a *userApi) CheckNickname(r *ghttp.Request) {
	var (
		data *model.UserApiCheckNickNameReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data.Nickname != "" && !service.User.CheckNickname(data.Nickname) {
		response.JsonExit(r, 1, "昵称已存在", false)
	}
	response.JsonExit(r, 0, "", true)
}

// @summary 获取用户详情信息
// @tags 用户服务
// @produce json
// @success 200 {object} model.User "执行结果"
// @router /user/profile [GET]
func (a *userApi) Profile(r *ghttp.Request) {
	response.JsonExit(r, 0, "", service.User.GetProfile(r.Context()))
}
