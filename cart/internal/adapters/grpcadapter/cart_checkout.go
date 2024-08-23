package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
)

func (g *GrpcAdapter) CartCheckout(ctx context.Context, in *grpc_cart.CartCheckoutRequest) (*grpc_cart.CartCheckoutResponse, error) {
	orderID, err := g.cartManager.CartCheckout(ctx, model.User{ID: in.UserID})
	if err != nil {
		return nil, err
	}
	return &grpc_cart.CartCheckoutResponse{OrderID: uint32(orderID)}, nil
}
