package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
)

func (g *GrpcAdapter) StockInfo(ctx context.Context, in *grpc_loms.StockInfoRequest) (*grpc_loms.StockInfoResponse, error) {
	stockInfo, err := g.lomsManager.StockInfo(ctx, model.SKU(in.Sku))
	if err != nil {
		return nil, err
	}
	return &grpc_loms.StockInfoResponse{
		Count: uint32(stockInfo),
	}, nil
}
