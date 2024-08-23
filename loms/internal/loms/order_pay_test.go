package loms

import (
	context "context"
	"database/sql"
	model "ecom/loms/internal/model"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderPay(t *testing.T) {
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
		expectedErr       error
		expectedOrderInfo model.OrderInfo
	}{
		{
			name:    "success",
			orderID: 1,
			orderRepoExpect: func() {
				orderInfoRepoReturn := model.OrderInfo{
					ID: 1,
					Items: []model.Item{
						{
							SKU:   1,
							Count: 1,
						},
					},
					StatusCode: model.OrderStatusAwaitingPayment,
				}
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(orderInfoRepoReturn, nil).Times(1)
				repoOrder.EXPECT().OrderPay(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(nil).Times(1)
			},
			expectedOrderInfo: model.OrderInfo{ID: 1},
		},
		{
			name:    "no such order",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Any()).Return(model.OrderInfo{}, model.ErrOrderIDNotFound).Times(1)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
		{
			name:    "order has wrong status",
			orderID: 1,
			orderRepoExpect: func() {
				orderInfoRepoReturn := model.OrderInfo{
					ID: 1,
					Items: []model.Item{
						{
							SKU:   1,
							Count: 1,
						},
					},
					StatusCode: model.OrderStatusFailed,
				}
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Any()).Return(orderInfoRepoReturn, nil).Times(1)
			},
			expectedErr: model.ErrOrderIsNotWaitingForPayment,
		},
		{
			name:    "some storage err",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Any()).Return(model.OrderInfo{}, sql.ErrNoRows).Times(1)
			},
			expectedErr: sql.ErrNoRows,
		},
		{
			name:    "storage orderPay err",
			orderID: 1,
			orderRepoExpect: func() {
				orderInfoRepoReturn := model.OrderInfo{
					ID: 1,
					Items: []model.Item{
						{
							SKU:   1,
							Count: 1,
						},
					},
					StatusCode: model.OrderStatusAwaitingPayment,
				}
				repoOrder.EXPECT().OrderInfoByOrderID(gomock.Any(), gomock.Any()).Return(orderInfoRepoReturn, nil).Times(1)
				repoOrder.EXPECT().OrderPay(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows).Times(1)
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.orderRepoExpect()
			ctx := context.Background()
			err := service.OrderPay(ctx, tc.orderID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
		})
	}
}
