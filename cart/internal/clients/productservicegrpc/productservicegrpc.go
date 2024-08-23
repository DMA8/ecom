package productservicegrpc

import (
	"context"
	"ecom/cart/internal/model"
	product "ecom/cart/pkg/api/productService"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type ProductServiceRatelim struct {
	cli product.ProductServiceClient
}

var (
	ErrProductService          = errors.New("product service error")
	ErrProductServiceRateLimit = errors.New("product service rate limit")
)

func New(cli product.ProductServiceClient) *ProductServiceRatelim {
	return &ProductServiceRatelim{
		cli: cli,
	}
}

func (p *ProductServiceRatelim) GetProduct(ctx context.Context, sku model.SKU) (model.ProductInfo, error) {
	const op = `ProductServiceCli.GetProduct`

	token, ok := ctx.Value(model.ProductServiceTokenKey).(string)
	if !ok {
		return model.ProductInfo{}, fmt.Errorf("%s err: %w", op, model.ErrNoTokenForProductService)
	}
	resp, err := p.cli.GetProduct(ctx, &product.GetProductRequest{
		Sku:   uint32(sku),
		Token: token,
	})
	if err != nil {
		return model.ProductInfo{}, fmt.Errorf("%s couldn't get product (sku=%d) %w", op, sku, err)
	}
	return model.ProductInfo{
		Name:  resp.Name,
		Price: resp.Price,
	}, nil
}

type RateLimInterceptor struct {
	requestCounter int64
	rps            int64
}

func NewRateLimInterceptor(rps int64) *RateLimInterceptor {
	return &RateLimInterceptor{
		rps: rps,
	}
}

func (r *RateLimInterceptor) RateLimiterInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	if atomic.AddInt64(&r.requestCounter, 1) > r.rps {
		return ErrProductServiceRateLimit
	}
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func (r *RateLimInterceptor) ResetCounterWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			atomic.StoreInt64(&r.requestCounter, 0)
		}
	}
}
