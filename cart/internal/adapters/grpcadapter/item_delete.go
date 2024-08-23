package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
)

func (g *GrpcAdapter) ItemDelete(ctx context.Context, in *grpc_cart.ItemDeleteRequest) (*grpc_cart.ItemDeleteResponse, error) {
	err := g.cartManager.ItemDelete(ctx, model.SKU(in.Sku), model.User{ID: in.UserID})
	if err != nil {
		return nil, err
	}
	return &grpc_cart.ItemDeleteResponse{}, nil
}
