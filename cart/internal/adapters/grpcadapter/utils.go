package grpcadapter

import (
	"ecom/cart/internal/model"
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
	case errors.Is(serviceError, model.ErrNoTokenForProductService):
		return status.Error(codes.Unauthenticated, "something wrong with product service token")
	case errors.Is(serviceError, model.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(serviceError, model.ErrCartIsEmpty):
		return status.Error(codes.FailedPrecondition, "cart is empty")
	case errors.Is(serviceError, model.ErrInvalidSKU):
		return status.Error(codes.InvalidArgument, "invalid sku")
	}
	return status.Error(codes.Internal, fmt.Sprintf("internal error: %s", serviceError.Error()))
}
