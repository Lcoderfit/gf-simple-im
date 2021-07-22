package main

import (
	_ "gf-simple-im/boot"
	_ "gf-simple-im/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
