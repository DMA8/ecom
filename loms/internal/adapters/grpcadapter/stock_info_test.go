package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestStockInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderManager := NewMockorderManager(ctrl)
	adapter := New(orderManager)

	testCases := []struct {
		name            string
		stockInfoReq    *grpc_loms.StockInfoRequest
		mockExpectation func()
		expectedErr     error
		expectedOrderID int64
	}{
		{
			name: "success",
			stockInfoReq: &grpc_loms.StockInfoRequest{
				Sku: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().StockInfo(gomock.Any(), model.SKU(1)).Return(int64(1), nil)
			},
			expectedOrderID: 1,
		},
		{
			name: "error",
			stockInfoReq: &grpc_loms.StockInfoRequest{
				Sku: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().StockInfo(gomock.Any(), model.SKU(1)).Return(int64(0), model.ErrOrderIDNotFound)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpectation()
			ctx := context.Background()
			_, err := adapter.StockInfo(ctx, tc.stockInfoReq)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
		})
	}
}
