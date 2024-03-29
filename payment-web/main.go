package main

import (
	"fmt"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-web/handler"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/utils"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry/consul"

	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
	"github.com/micro/go-micro/web"
	"net/http"
)

var (
	log = logzap.GetLogger()
)

var (
	appName = "payment_web"
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
	// 创建新服务
	service := web.NewService(
		web.Name(cfg.Name),
		web.Version(cfg.Version),
		web.Registry(micReg),
	)

	// 初始化服务
	if err := service.Init(
		web.Action(
			func(c *cli.Context) {
				// 初始化handler
				handler.Init()
			}),
	); err != nil {
		log.Fatal(err.Error())
	}

	// 新建订单接口
	authHandler := http.HandlerFunc(handler.PayOrder)
	service.Handle("/payment/pay-order", utils.AuthWrapper(authHandler))

	// 运行服务
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
