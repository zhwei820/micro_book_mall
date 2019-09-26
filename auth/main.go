package main

import (
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"go.uber.org/zap"
	"time"

	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/handler"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/model"
	s "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/proto/auth"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	z "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
)

var (
	log     = z.GetLogger()
	appName = "auth_srv"
	cfg     = &authCfg{}
)

type authCfg struct {
	common.AppCfg
}

func main() {
	// 初始化配置、数据库等信息
	basic.InitAppCfg(appName)
	InitCfg()

	// 使用consul注册
	micReg := consul.NewRegistry(basic.RegistryOptions)

	// 新建服务
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Registry(micReg),
		micro.Version(cfg.Version),

		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
	)

	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 初始化handler
			model.Init()
			// 初始化handler
			handler.Init()
		}),
	)

	// 注册服务
	s.RegisterServiceHandler(service.Server(), new(handler.Service))

	// 启动服务
	if err := service.Run(); err != nil {
		log.Error("[main] error")
		panic(err)
	}
}

func InitCfg() {

	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}

	log.Info("[InitAppCfg] 配置", zap.Any("cfg", cfg))

	return
}
