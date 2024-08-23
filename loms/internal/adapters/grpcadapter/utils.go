package grpcadapter

import (
	"ecom/loms/internal/model"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertError(serviceError error) error {
	if serviceError == nil {
		return nil
	}
	switch {
	case errors.Is(serviceError, model.ErrOrderIDNotFound):
		return status.Error(codes.NotFound, "order not found")
	case errors.Is(serviceError, model.ErrNotEnoughStock):
		return status.Error(codes.FailedPrecondition, "not enough stock")
	case errors.Is(serviceError, model.ErrSKUNotValid):
		return status.Error(codes.InvalidArgument, "sku not valid")
	case errors.Is(serviceError, model.ErrOrderIsNotWaitingForPayment):
		return status.Error(codes.FailedPrecondition, "order is not waiting for payment")
	case errors.Is(serviceError, model.ErrSKUNotFound):
		return status.Error(codes.NotFound, "sku not found")
	}
	return status.Error(codes.Internal, fmt.Sprintf("internal error: %s", serviceError.Error()))
}
