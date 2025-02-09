package service

import (
	"backend/application/product/internal/biz"
	"context"

	pb "backend/api/product/v1"
)

type ProductCatalogServiceService struct {
	pb.UnimplementedProductCatalogServiceServer
	pc *biz.ProductUsecase
}

func (s *ProductCatalogServiceService) UpdateProduct(ctx context.Context, product *pb.Product) (*pb.ProductReply, error) {
	// TODO implement me
	panic("implement me")
}

func (s *ProductCatalogServiceService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	products, err := s.ps.SearchProducts(ctx, &biz.SearchProductsReq{Query: req.GetQuery()})
	if err != nil {
		return nil, err
	}
	pbProduct := make([]*pb.Product, len(products.Result))
	for i, product := range products.Result {
		pbProduct[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  product.Categories,
		}
	}
	return &pb.SearchProductsResp{
		Results: pbProduct,
	}, nil
}

func (s *ProductCatalogServiceService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductReply, error) {
	p, err := s.ps.CreateProduct(ctx, &biz.CreateProductRequest{
		Owner:       req.Owner,
		Username:    req.Username,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{
		Product: &pb.Product{
			Id:          p.Product.Id,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Picture:     p.Product.Picture,
			Price:       p.Product.Price,
			Categories:  p.Product.Categories,
		},
	}, nil
}

func (s *ProductCatalogServiceService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	list, err := s.ps.ListProducts(ctx, &biz.ListProductsReq{
		Page:         uint(req.Page),
		PageSize:     uint(req.PageSize),
		CategoryName: req.CategoryName,
	})
	if err != nil {
		return nil, err
	}
	pbProduct := make([]*pb.Product, len(list.Product))
	for i, product := range list.Product {
		pbProduct[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  product.Categories,
		}
	}
	return &pb.ListProductsResp{
		Products: pbProduct,
	}, nil
}

func (s *ProductCatalogServiceService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.ProductReply, error) {
	product, err := s.ps.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{Product: &pb.Product{
		Id:          product.Product.Id,
		Name:        product.Product.Name,
		Description: product.Product.Description,
		Picture:     product.Product.Picture,
		Price:       product.Product.Price,
		Categories:  product.Product.Categories,
	}}, nil
}

func NewServiceProductCatalogServiceService(ps *biz.ProductUsecase) *ProductCatalogServiceService {
	return &ProductCatalogServiceService{ps: ps}
}
