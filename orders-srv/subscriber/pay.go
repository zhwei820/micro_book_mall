package subscriber

import (
	"context"
	"fmt"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/orders-srv/model/orders"
	payS "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-srv/proto/payment"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log = logzap.GetLogger()
)

var (
	ordersService orders.Service
)

// Init 初始化handler
func Init() {
	ordersService, _ = orders.GetService()
}

// PayOrder 订单支付消息
func PayOrder(ctx context.Context, event *payS.PayEvent) (err error) {
	log.Info(fmt.Sprintf("[PayOrder] 收到支付订单通知，%d，%d", event.OrderId, event.State))
	err = ordersService.UpdateOrderState(event.OrderId, int(event.State))
	if err != nil {
		log.Info(fmt.Sprintf("[PayOrder] 收到支付订单通知，更新状态异常，%s", err))
		return
	}
	return
}
