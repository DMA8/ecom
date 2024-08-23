package cart

import (
	"context"
	"ecom/cart/internal/model"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
)

func fakeApiCall(ctx context.Context, sku model.SKU) (model.ProductInfo, error) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if rand.Intn(100) > 110 {
		return model.ProductInfo{}, fmt.Errorf("fakeApiCall error")
	}
	return model.ProductInfo{}, nil
}

func TestLeakErrGroupGetProductForItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer leaktest.Check(t)()
	itemsTest := make([]model.Item, 100)
	getProductForItemsErrGroups(ctx, itemsTest, fakeApiCall)
}

func TestLeaksAsyncGetProductForItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer leaktest.Check(t)()
	itemsTest := make([]model.Item, 100)
	_, err := getProductForItemsAsync(ctx, itemsTest, fakeApiCall)
	if err != nil {
		log.Println(err)
	}
}
