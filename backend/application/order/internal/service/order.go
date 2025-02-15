package service

import (
	v1 "backend/api/cart/v1"
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
	oc         *biz.OrderUsecase
	CartClient v1.CartServiceClient // 引入购物车服务的客户端
}

func NewOrderServiceService(cartClient v1.CartServiceClient) *OrderServiceService {
	return &OrderServiceService{
		CartClient: cartClient,
	}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq, userID string) (*pb.PlaceOrderResp, error) {

	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	//调用购物车服务
	CartResp, err := s.CartClient.GetCart(ctx, &v1.GetCartReq{
		Owner: userID,
		//Name:  req.Name,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "获取购物车失败")
	}

	UserId, err := strconv.ParseUint(payload.ID, 10, 32)
	if err != nil {
		// 处理转换错误，例如返回错误信息给调用者
		return nil, status.Error(codes.Internal, "用户ID转换失败")
	}

	var items []biz.Item
	for _, cartItem := range CartResp.Cart.Items {
		items = append(items, biz.Item{
			Id:       int32(cartItem.ProductId),
			Quantity: cartItem.Quantity,
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
			OrderId: string(result.Order.OrderId),
		},
	}, nil

}

func (s *OrderServiceService) ListOrders(ctx context.Context, req *pb.ListOrderReq) (*pb.ListOrderResp, error) {

	bizResp, err := s.oc.ListOrders(ctx, &biz.ListOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range bizResp.Orders {
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
