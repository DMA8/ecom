//go:generate mockgen -source $GOFILE -destination mock_contract.go -package ${GOPACKAGE}
package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"

	"google.golang.org/grpc"
)

type orderManager interface {
	StockInfo(ctx context.Context, sku model.SKU) (int64, error)
	OrderCreate(ctx context.Context, order model.OrderCreate, user model.User) (model.OrderID, error)
	OrderInfoByOrderID(ctx context.Context, orderID model.OrderID) (model.OrderInfo, error)
	OrderPay(ctx context.Context, orderID model.OrderID) error
	OrderCancel(ctx context.Context, orderID model.OrderID) error
}

type GrpcAdapter struct {
	lomsManager orderManager
	grpc_loms.UnimplementedLOMSServiceServer
}

func New(lm orderManager) *GrpcAdapter {
	return &GrpcAdapter{
		lomsManager: lm,
	}
}

func (g *GrpcAdapter) Server() *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			logInterceptor,
			validateInterceptor,
			errorConverter,
			productServiceTokenValidatorInterceptor,
		),
	)
	grpc_loms.RegisterLOMSServiceServer(s, g)
	return s
}
