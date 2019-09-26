package handler

import (
	"context"
	hystrix_go "github.com/afex/hystrix-go/hystrix"
	auth "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/proto/auth"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
	us "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/user-srv/proto/user"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"strconv"
)

var (
	log = logzap.GetLogger()
)

var (
	serviceClient us.UserService
	authClient    auth.Service
)

func Init() {
	hystrix_go.DefaultVolumeThreshold = 1
	hystrix_go.DefaultErrorPercentThreshold = 1
	cl := hystrix.NewClientWrapper()(client.DefaultClient)
	cl.Init(
		client.Retries(3),
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			log.Info(ctx, req.Method()+" "+strconv.Itoa(retryCount)+" client retry")
			return true, nil
		}),
	)
	serviceClient = us.NewUserService("mu.micro.book.srv.user", cl)
	authClient = auth.NewService("mu.micro.book.srv.auth", cl)
}
