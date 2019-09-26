package main

import (
	"fmt"
	"time"

	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/inventory-srv/handler"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/inventory-srv/model"
	proto "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/inventory-srv/proto/inventory"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
)

var (
	log = logzap.GetLogger()
)

var (
	appName = "inv_srv"
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
		micro.Registry(micReg),
		micro.Version(cfg.Version),

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
	proto.RegisterInventoryHandler(service.Server(), new(handler.Service))

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
