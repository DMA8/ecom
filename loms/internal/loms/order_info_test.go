package loms

import (
	context "context"
	model "ecom/loms/internal/model"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoStock := NewMockStockRepository(ctrl)
	repoOrder := NewMockOrderRepository(ctrl)
	producer := NewMockProducer(ctrl)
	service := New(repoOrder, repoStock, producer)

	testCases := []struct {
		name              string
		orderID           model.OrderID
		orderRepoExpect   func()
		repoOrderExpect   func()
		expectedErr       error
		expectedOrderInfo model.OrderInfo
	}{
		{
			name:    "success",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(model.OrderInfo{ID: 1, Items: []model.Item{{SKU: 1, Count: 1}}}, nil).Times(1)
			},
			expectedOrderInfo: model.OrderInfo{ID: 1},
		},
		{
			name:    "error",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(model.OrderInfo{}, nil).Times(1)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.orderRepoExpect()
			ctx := context.Background()
			orderInfo, err := service.OrderInfoByOrderID(ctx, tc.orderID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
			if orderInfo.ID != tc.expectedOrderInfo.ID {
				t.Errorf("name: %s. id expected %v, got %v", tc.name, tc.expectedOrderInfo.ID, orderInfo)
			}
		})
	}
}
