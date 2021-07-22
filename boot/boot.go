package boot

import (
	_ "gf-simple-im/packed"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

// 用于应用初始化
func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
}
