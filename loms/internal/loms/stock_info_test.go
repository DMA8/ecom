package loms

import (
	"context"
	model "ecom/loms/internal/model"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestStockInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoStock := NewMockStockRepository(ctrl)
	repoOrder := NewMockOrderRepository(ctrl)
	producer := NewMockProducer(ctrl)
	service := New(repoOrder, repoStock, producer)

	testCases := []struct {
		name             string
		sku              model.SKU
		stockRepoExpect  func()
		expectedQuantity int64
		expectedErr      error
	}{
		{
			name: "success",
			sku:  model.SKU(1),
			stockRepoExpect: func() {
				repoStock.EXPECT().StockQuantity(gomock.Any(), gomock.Eq(model.SKU(1))).Return(int64(10), nil).Times(1)
			},
			expectedQuantity: 10,
		},
		{
			name: "error",
			sku:  model.SKU(1),
			stockRepoExpect: func() {
				repoStock.EXPECT().StockQuantity(gomock.Any(), gomock.Eq(model.SKU(1))).Return(int64(0), model.ErrSKUNotFound).Times(1)
			},
			expectedQuantity: 0,
			expectedErr:      model.ErrSKUNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.stockRepoExpect()
			ctx := context.Background()
			q, err := service.StockInfo(ctx, tc.sku)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("err expected %v, got %v", tc.expectedErr, err)
			}
			if q != tc.expectedQuantity {
				t.Errorf("quantity expected %v, got %v", tc.expectedQuantity, q)
			}
		})
	}
}
