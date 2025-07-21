package main

import (
	"flag"
	"fmt"

	"short/internal/config"
	"short/internal/handler"
	"short/internal/pkg/base62"
	"short/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/short-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c) // 加载配置
	fmt.Printf("load config: %#v\n", c)

	base62.MustInit(c.Base62String) // base62模块初始化, 从配置加载

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
