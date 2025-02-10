package service

import (
	pb "backend/api/order/v1"
	"backend/application/order/internal/biz"
	"backend/application/order/pkg/convert"
	"backend/application/order/pkg/token"
	"context"
)

type OrderServiceService struct {
	pb.UnimplementedOrderServiceServer
	oc *biz.OrderUsecase
}

func NewOrderServiceService() *OrderServiceService {
	return &OrderServiceService{}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderResp, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	var items []biz.Item
	for _, item := range req.Items {
		items = append(items, biz.Item{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}

	result, err := s.oc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		Name:         payload.Name,
		UserCurrency: req.UserCurrency,
		Address: biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			ZipCode:       req.Address.ZipCode,
		},
		Items: items,
		Email: req.Email,
		Owner: payload.Owner,
	})
	if err != nil {
		return nil, err
	}

	return &pb.PlaceOrderResp{
		Order: &pb.OrderResult{
			OrderId: result.Order.OrderId,
		},
	}, nil
}
func (s *OrderServiceService) ListOrder(ctx context.Context, req *pb.ListOrderReq) (*pb.ListOrderResp, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	orders, err := s.oc.ListOrder(ctx, &biz.ListOrderReq{
		UserId: payload.ID,
	})
	if err != nil {
		return nil, err
	}

	// 将 biz.Order 转换为 pb.Order
	var pbOrders []*pb.Order
	for _, order := range orders.Orders {
		pbOrders = append(pbOrders, &pb.Order{
			OrderId:      order.OrderId,
			UserId:       order.UserId,
			UserCurrency: order.UserCurrency,
			Email:        order.Email,
			CreatedAt:    order.CreatedAt,
			OrderItems:   convert.ToPbOrderItems(order.OrderItems),
			Address: &pb.Address{
				StreetAddress: order.Address.StreetAddress,
				City:          order.Address.City,
				State:         order.Address.State,
				ZipCode:       order.Address.ZipCode,
			},
		})
	}

	// 返回响应
	return &pb.ListOrderResp{
		Orders: pbOrders,
	}, nil
}

// func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidReq) (*pb.MarkOrderPaidResp, error) {
// 	return &pb.MarkOrderPaidResp{}, nil
// }
