package service

import (
	"context"
	"gf-simple-im/app/model"
	"github.com/gogf/gf/net/ghttp"
)

// 上下文管理服务
var Context = contextService{}

type contextService struct{}

// 初始化自定义上下文变量指针到gf上下文变量中，方便后续修改
func (s *contextService) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(model.ContextKey, customCtx)
}

// 获取上下文变量，如果没有设置，那么返回nil
func (s *contextService) Get(ctx context.Context) *model.Context {
	value := ctx.Value(model.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// 将ContextUser指针赋给自定义Context的User字段
func (s *contextService) SetUser(ctx context.Context, ctxUser *model.ContextUser) {
	s.Get(ctx).User = ctxUser
}
