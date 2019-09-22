package env_config

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/env"
)

var EnvConfig = struct {
	ConfigAddress string
}{}

func init() {
	src := env.NewSource(
		// optionally specify prefix
		//env.WithPrefix("MICRO"),
	)
	// Create new config
	conf := config.NewConfig()

	// Load env source
	conf.Load(src)
	conf.Scan(&EnvConfig)
	fmt.Printf("%v", EnvConfig)
}
