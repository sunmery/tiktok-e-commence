package data

import (
	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/application/order/pkg/convert"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	// 创建订单
	orderID, err := o.data.db.CreateOrder(ctx, models.CreateOrderParams{
		Owner:         req.Owner,
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

// func (o *orderRepo) ListOrders(ctx context.Context, req *biz.ListOrderReq) ([]*biz.ListOrderResp, error) {
// 	// 提取 payload
// 	payload, err := token.ExtractPayload(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 调用数据库查询
// 	dbOrders, err := o.data.db.ListOrders(ctx, ListOrdersParams{
// 		Owner: payload.Owner,
// 		Name:  req.Name,
// 	})
// }

func NewOrderrRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
