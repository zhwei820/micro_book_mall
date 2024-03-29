package handler

import (
	"context"
	"fmt"

	inv "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/inventory-srv/model/inventory"
	proto "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/inventory-srv/proto/inventory"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log = logzap.GetLogger()
)

var (
	invService inv.Service
)

type Service struct {
}

// Init 初始化handler
func Init() {
	invService, _ = inv.GetService()
}

// Sell 库存销存
func (e *Service) Sell(ctx context.Context, req *proto.Request, rsp *proto.Response) (err error) {
	id, err := invService.Sell(req.BookId, req.UserId)
	if err != nil {
		log.Info(fmt.Sprintf("[Sell] 销存失败，bookId：%d，userId: %d，%s", req.BookId, req.UserId, err))
		rsp.Success = false
		return
	}

	rsp.InvH = &proto.InvHistory{
		Id: id,
	}

	rsp.Success = true
	return nil
}

// Confirm 库存销存 确认
func (e *Service) Confirm(ctx context.Context, req *proto.Request, rsp *proto.Response) (err error) {
	err = invService.Confirm(req.HistoryId, int(req.HistoryState))
	if err != nil {
		log.Info(fmt.Sprintf("[Confirm] 确认销存失败，%s", err))
		rsp.Success = false
		return
	}

	rsp.Success = true
	return nil
}
