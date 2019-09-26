package basic

import (
	"fmt"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/utils/env-config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/config/source/grpc"
)

var (
	pluginFuncs []func()
)

type Options struct {
	EnableDB    bool
	EnableRedis bool
	cfgOps      []config.Option
}

type Option func(o *Options)

func Init(opts ...config.Option) {
	// 初始化配置
	config.Init(opts...)

	// 加载依赖配置的插件
	for _, f := range pluginFuncs {
		f()
	}
}

func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}

func RegistryOptions(ops *registry.Options) {
	consulCfg := &common.Consul{}
	err := config.C().App("consul", consulCfg)
	if err != nil {
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", consulCfg.Host, consulCfg.Port)}
}

func InitAppCfg(appName string) {
	source := grpc.NewSource(
		grpc.WithAddress(env_config.EnvConfig.ConfigAddress),
		grpc.WithPath("micro"),
	)

	Init(
		config.WithSource(source),
		config.WithApp(appName),
	)

}
