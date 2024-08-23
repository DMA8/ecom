package grpcadapter

import (
	"context"
	"ecom/cart/internal/model"
	grpc_cart "ecom/cart/pkg/api/cart/v1"
	"fmt"
)

func (g *GrpcAdapter) ItemAdd(ctx context.Context, in *grpc_cart.ItemAddRequest) (*grpc_cart.ItemAddResponse, error) {
	itemAddModel, err := itemAddDto(in)
	if err != nil {
		return nil, err
	}

	err = g.cartManager.AddItem(ctx, itemAddModel, model.User{ID: in.UserID})
	if err != nil {
		return nil, err
	}
	return &grpc_cart.ItemAddResponse{}, nil
}

func itemAddDto(in *grpc_cart.ItemAddRequest) (model.Item, error) {
	if in.ItemToAdd == nil {
		return model.Item{}, fmt.Errorf("ItemAdd: %w", ErrUnexpectedEmptyInputGrpcMethod)
	}
	return model.Item{
		SKU:   model.SKU(in.ItemToAdd.Sku),
		Count: uint16(in.ItemToAdd.Count),
	}, nil
}
