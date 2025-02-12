package service

import (
	pb "backend/api/order/v1"
	"backend/application/order/internal/biz"
	"backend/application/order/pkg/convert"
	"backend/application/order/pkg/token"
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	UserId, err := strconv.ParseUint(payload.ID, 10, 32)
	if err != nil {
		// 处理转换错误，例如返回错误信息给调用者
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
		UserId:       uint32(UserId),
		UserCurrency: req.UserCurrency,
		Address: biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			ZipCode:       req.Address.ZipCode,
		},
		Items: items,
		Email: req.Email,
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

	bizResp, err := s.oc.ListOrder(ctx, &biz.ListOrderReq{
		UserId: uint32(req.UserId),
	})
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range bizResp.Orders {
		pbOrders = append(pbOrders, &pb.Order{
			Status:       order.Status,
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

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidReq) (*pb.MarkOrderPaidResp, error) {

	log.Printf("MarkOrderPaid called with OrderId: %s, UserId: %d", req.OrderId, req.UserId)

	if req.OrderId == "" {
		log.Println("OrderId is empty")
		return nil, status.Error(codes.InvalidArgument, "订单ID不能为空")
	}

	_, err := s.oc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  uint32(req.UserId),
		OrderId: req.OrderId,
	})

	if err != nil {
		log.Printf("MarkOrderPaid failed: %v", err)
		return nil, status.Error(codes.Internal, "订单支付失败")
	}

	log.Println("MarkOrderPaid succeeded")

	return &pb.MarkOrderPaidResp{}, nil
}
