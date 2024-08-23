package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
	"errors"

	"google.golang.org/grpc"
)

var (
	ErrUnexpectedEmptyInputGrpcMethod = errors.New("unexpected empty input grpc method")
)

type CartManager interface {
	AddItem(ctx context.Context, item model.Item, user model.User) error
	ItemDelete(ctx context.Context, sku model.SKU, user model.User) error
	CartClear(ctx context.Context, user model.User) error
	CartList(ctx context.Context, user model.User) (model.UserCart, error)
	CartCheckout(ctx context.Context, user model.User) (model.OrderID, error)
}

type GrpcAdapter struct {
	cartManager CartManager
	grpc_cart.UnimplementedCartServiceServer
}

func New(cm CartManager) *GrpcAdapter {
	return &GrpcAdapter{
		cartManager: cm,
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
	grpc_cart.RegisterCartServiceServer(s, g)
	return s
}
