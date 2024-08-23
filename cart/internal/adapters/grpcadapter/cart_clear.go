package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
)

func (g *GrpcAdapter) CartClear(ctx context.Context, in *grpc_cart.CartClearRequest) (*grpc_cart.CartClearResponse, error) {
	err := g.cartManager.CartClear(ctx, model.User{ID: in.UserID})
	if err != nil {
		return nil, err
	}
	return &grpc_cart.CartClearResponse{}, nil
}
