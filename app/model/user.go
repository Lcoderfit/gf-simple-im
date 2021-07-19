package model

// 注册请求参数，用于前后端交互参数格式约定
type UserApiSignUpReq struct {
	Passport  string `v:"required|length:6,16#账号不能为空|账号长度应当在:min到:max之间"`
	Password  string `v:"required|length:6,16#密码不能为空|密码长度应当在:min到:max之间"`
	Password2 string `v:"required|length:6,16|same:Password#密码长度不能为空|密码长度应当在:min到max之间|两次密码输入不相等"`
	Nickname  string
}

// 登录请求参数
type UserApiSignInReq struct {
	Passport string `v:"required#账号不能为空"`
	Password string `v:"required#密码不能为空"`
}

// 账号唯一性检测请求参数
type UserApiCheckPassportReq struct {
	Passport string `v:"required#账号不能为空"`
}

//// User is the golang structure for table user.
//type User internal.User

// 昵称唯一性检测请求参数
type UserApiCheckNickNameReq struct {
	Nickname string `v:"required#昵称不能为空"`
}

// 注册输入参数
type UserServiceSignUpReq struct {
	Passport string
	Password string
	Nickname string
}
