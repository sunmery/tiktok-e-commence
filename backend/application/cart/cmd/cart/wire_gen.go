// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/cart/internal/biz"
	"backend/application/cart/internal/conf"
	"backend/application/cart/internal/data"
	"backend/application/cart/internal/server"
	"backend/application/cart/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, consul *conf.Consul, observability *conf.Observability, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	dataData, cleanup, err := data.NewData(pool, client, logger)
	if err != nil {
		return nil, nil, err
	}
	cartRepo := data.NewCartRepo(dataData, logger)
	cartUsecase := biz.NewCartUsecase(cartRepo, logger)
	cartServiceService := service.NewCartServiceService(cartUsecase)
	grpcServer := server.NewGRPCServer(cartServiceService, confServer, observability, logger)
	httpServer := server.NewHTTPServer(confServer, cartServiceService, auth, observability, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
