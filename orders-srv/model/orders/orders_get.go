package orders

import (
	"fmt"
	proto "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/orders-srv/proto/orders"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/db"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log = logzap.GetLogger()
)

// GetOrder 获取订单
func (s *service) GetOrder(orderId int64) (order *proto.Order, err error) {
	order = &proto.Order{}

	// 获取数据库
	o := db.GetDB()
	// 查询
	err = o.QueryRow("SELECT id, user_id, book_id, inv_his_id, state FROM orders WHERE id = ?", orderId).Scan(
		&order.Id, &order.UserId, &order.BookId, &order.InvHistoryId, &order.State)
	if err != nil {
		log.Info(fmt.Sprintf("[GetOrder] 查询数据失败，err：%s", err))
		return
	}

	return
}
