package main

import (
	"github.com/opentracing/opentracing-go"
	"log"

	tracer "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/tracer/jaeger"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/tracer/opentracing/stdhttp"
	"github.com/micro/go-plugins/micro/cors"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(cors.NewPlugin())

	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
}

const name = "API gateway"

func main() {
	stdhttp.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	cmd.Init()
}
