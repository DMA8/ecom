package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
)

func (g *GrpcAdapter) OrderCancel(ctx context.Context, in *grpc_loms.OrderCancelRequest) (*grpc_loms.OrderCancelResponse, error) {
	err := g.lomsManager.OrderCancel(ctx, model.OrderID(in.OrderId))
	if err != nil {
		return nil, err
	}
	return &grpc_loms.OrderCancelResponse{}, nil
}
