package main

import (
	"github.com/micro/go-plugins/micro/cors"
	"github.com/micro/go-plugins/micro/trace/uuid"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(cors.NewPlugin())
	plugin.Register(uuid.New())
}

func main() {
	cmd.Init()
}
