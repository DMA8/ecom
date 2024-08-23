package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ProductServiceTokenKey = "x-product-service-token"
)

type Validatable interface {
	ValidateAll() error
}

func logInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	ts := time.Now()
	resp, err := handler(ctx, req)
	slog.Info("grpc call", "method", info.FullMethod, "duration", time.Since(ts), "err", err)
	return resp, err
}

func validateInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	if reqValidatable, ok := req.(Validatable); ok {
		if err := reqValidatable.ValidateAll(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

func recoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	defer func() {
		if r := recover(); r != nil {
			_ = status.Errorf(codes.Internal, "Panic recovered: %v", r)
			slog.Error("Panic recovered: %v\n", r)
		}
	}()

	return handler(ctx, req)
}

func productServiceTokenValidatorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md.Get(ProductServiceTokenKey)
		if len(val) > 0 {
			ctx = context.WithValue(ctx, model.ProductServiceTokenKey, val[0])
		}
	}
	return handler(ctx, req)
}

func errorConverter(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	resp, err := handler(ctx, req)
	transportErr := convertError(err)
	return resp, transportErr
}
