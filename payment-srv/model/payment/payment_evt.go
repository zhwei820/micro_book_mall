package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/payment-srv/proto/payment"

	"time"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log = logzap.GetLogger()
)

// sendPayDoneEvt 发送支付事件
func (s *service) sendPayDoneEvt(orderId int64, state int32) {
	// 构建事件
	ev := &proto.PayEvent{
		Id:       uuid.New().String(),
		SentTime: time.Now().Unix(),
		OrderId:  orderId,
		State:    state,
	}

	log.Info(fmt.Sprintf("[sendPayDoneEvt] 发送支付事件，%+v\n", ev))

	// 广播
	if err := payPublisher.Publish(context.Background(), ev); err != nil {
		log.Info(fmt.Sprintf("[sendPayDoneEvt] 异常: %v", err))
	}
}
