package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
)

func (g *GrpcAdapter) CartList(ctx context.Context, in *grpc_cart.CartListRequest) (*grpc_cart.CartListResponse, error) {
	cart, err := g.cartManager.CartList(ctx, model.User{ID: in.UserID})
	if err != nil {
		return nil, err
	}

	response := prepareResponseCartList(cart)
	return &response, nil
}

func prepareResponseCartList(cart model.UserCart) grpc_cart.CartListResponse {

	itemsOut := make([]*grpc_cart.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		itemsOut = append(itemsOut, &grpc_cart.CartItem{
			Sku:   uint32(item.Sku),
			Name:  item.Name,
			Price: item.Price,
			Count: uint32(item.Count),
		})
	}

	return grpc_cart.CartListResponse{
		TotalPrice: cart.TotalPrice,
		CartItems:  itemsOut,
	}
}
