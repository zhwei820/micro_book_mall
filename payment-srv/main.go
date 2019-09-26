package main

import (
	"fmt"
	"time"

	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-srv/handler"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-srv/model"
	s "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-srv/proto/payment"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"

	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log = logzap.GetLogger()
)

var (
	appName = "payment_srv"
	cfg     = &appCfg{}
)

type appCfg struct {
	common.AppCfg
}

func main() {
	// 初始化配置、数据库等信息
	basic.InitAppCfg(appName)
	initCfg()

	// 使用consul注册
	micReg := consul.NewRegistry(basic.RegistryOptions)

	// 新建服务
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Registry(micReg),

		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
	)

	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 初始化模型层
			model.Init()
			// 初始化handler
			handler.Init()
		}),
	)

	// 注册服务
	s.RegisterPaymentHandler(service.Server(), new(handler.Service))

	// 启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func initCfg() {

	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("[initCfg] 配置，cfg：%v", cfg))

	return
}
