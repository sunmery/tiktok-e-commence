package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// CreditCards is a CreditCards model.
type CreditCards struct {
	Id              uint32 `json:"id"`
	Owner           string `json:"owner"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	Cvv             string `json:"cvv"`
	ExpirationYear  string `json:"expiration_year"`
	ExpirationMonth string `json:"expiration_month"`
}

type GetCreditCardsRequest struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Number string `json:"number"`
}
type CreditCardsRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreditCardsReply struct {
	Message string `json:"message"`
	Code    int32  `json:"code"`
}

type DeleteCreditCardsRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Id    uint32 `json:"id"`
}

// CreditCardsRepo is a Greater repo.
type CreditCardsRepo interface {
	CreateCreditCard(ctx context.Context, req *CreditCards) (*CreditCardsReply, error)
	UpdateCreditCard(ctx context.Context, req *CreditCards) (*CreditCardsReply, error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardsRequest) (*CreditCardsReply, error)
	GetCreditCard(ctx context.Context, req *GetCreditCardsRequest) (*CreditCards, error)
	SearchCreditCards(ctx context.Context, req *GetCreditCardsRequest) ([]*CreditCards, error)
	ListCreditCards(ctx context.Context, req *CreditCardsRequest) ([]*CreditCards, error)
}

// CreditCardsUsecase is a CreditCards usecase.
type CreditCardsUsecase struct {
	repo CreditCardsRepo
	log  *log.Helper
}

func (cc *UserUsecase) CreateCreditCards(ctx context.Context, req *CreditCards) (*CreditCardsReply, error) {
	cc.log.WithContext(ctx).Infof("CreateCreditCards: %+v\n", req)
	return cc.repo.CreateCreditCard(ctx, req)
}
func (cc *UserUsecase) UpdateCreditCards(ctx context.Context, req *CreditCards) (*CreditCardsReply, error) {
	cc.log.WithContext(ctx).Infof("UpdateCreditCards: %+v\n", req)
	return cc.repo.UpdateCreditCard(ctx, req)
}
func (cc *UserUsecase) DeleteCreditCards(ctx context.Context, req *DeleteCreditCardsRequest) (*CreditCardsReply, error) {
	cc.log.WithContext(ctx).Infof("DeleteCreditCards: %+v\n", req)
	return cc.repo.DeleteCreditCard(ctx, req)
}
func (cc *UserUsecase) GetCreditCard(ctx context.Context, req *GetCreditCardsRequest) (*CreditCards, error) {
	cc.log.WithContext(ctx).Infof("GetCreditCards: %+v\n", req)
	return cc.repo.GetCreditCard(ctx, req)
}
func (cc *UserUsecase) SearchCreditCards(ctx context.Context, req *GetCreditCardsRequest) ([]*CreditCards, error) {
	cc.log.WithContext(ctx).Infof("GetCreditCards: %+v\n", req)
	return cc.repo.SearchCreditCards(ctx, req)
}
func (cc *UserUsecase) ListCreditCards(ctx context.Context, req *CreditCardsRequest) ([]*CreditCards, error) {
	cc.log.WithContext(ctx).Infof("ListCreditCards: %+v\n", req)
	return cc.repo.ListCreditCards(ctx, req)
}
