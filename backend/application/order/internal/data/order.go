package data

import (
	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/application/order/pkg/convert"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
)

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	// 创建订单
	orderID, err := o.data.db.CreateOrder(ctx, models.CreateOrderParams{
		Owner:         fmt.Sprintf("%d", req.UserId),
		Name:          req.Name,
		Email:         req.Email,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		ZipCode:       req.Address.ZipCode,
		Currency:      req.UserCurrency,
	})
	if err != nil {
		return nil, err
	}

	// 创建订单商品
	for _, item := range req.Items {
		price, err := convert.Float32ToNumeric(item.Price)
		if err != nil {
			return nil, err
		}

		_, err = o.data.db.CreateOrderItems(ctx, models.CreateOrderItemsParams{
			OrderID:   orderID,
			ProductID: item.Id,
			Name:      item.Name,
			Price:     price,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return nil, err
		}
	}

	// 返回响应
	return &biz.PlaceOrderResp{
		Order: biz.OrderResult{
			OrderId: orderID,
		},
	}, nil
}

func (o *orderRepo) ListOrders(ctx context.Context, req *biz.ListOrderReq) ([]*biz.ListOrderResp, error) {

	// 从数据库获取订单数据
	dbOrders, err := o.data.db.ListOrders(ctx, models.ListOrdersParams{
		Owner: fmt.Sprintf("%d", req.UserId),
	})
	if err != nil {
		return nil, err
	}

	// 将数据库模型转换为业务模型
	var orderSummaries []biz.OrderSummary
	for _, dbOrder := range dbOrders {
		orderSummaries = append(orderSummaries, biz.OrderSummary{
			OrderId:   dbOrder.ID,
			CreatedAt: dbOrder.CreatedAt.Unix(),
			Address: biz.Address{ // 假设 dbOrder 包含了地址信息
				StreetAddress: dbOrder.StreetAddress,
				City:          dbOrder.City,
				State:         dbOrder.State,
				Country:       dbOrder.Country,
				ZipCode:       dbOrder.ZipCode,
			},
			Status:       dbOrder.Status,
			UserCurrency: dbOrder.Currency,
			Email:        dbOrder.Email,
		})
	}

	// 返回包含订单列表的响应
	return []*biz.ListOrderResp{{
		Orders: orderSummaries,
	}}, nil
}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// 1. 参数校验
	if req.OrderId == "" {
		return nil, errors.New("订单ID不能为空")
	}
	if req.UserId == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	orderId, err := strconv.ParseInt(req.OrderId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID格式: %w", err)
	}

	// 调用 casdoor 接口标记订单为已支付
	_, err = o.data.db.MarkOrderPaid(ctx, models.MarkOrderPaidParams{
		ID:    int32(orderId),
		Owner: fmt.Sprintf("%d", req.UserId),
		Name:  "test",
	})

	if err != nil {
		// 根据错误类型返回具体的错误信息
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("订单不存在")
		}
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return &biz.MarkOrderPaidResp{
		Success: true,
	}, nil

}

func NewOrderrRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
