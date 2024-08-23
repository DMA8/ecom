package cart

import (
	"context"
	"ecom/cart/internal/model"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

func getProductForItemsAsync(ctx context.Context, items []model.Item, apiCall func(ctx context.Context, sku model.SKU) (model.ProductInfo, error)) (map[model.SKU]ProductInfo, error) {
	var errApiCall error
	const op = `cart.getProductForItemsAsync`

	if len(items) < 1 {
		return make(map[model.SKU]ProductInfo), nil
	}

	skuToGetProductChan := make(chan model.SKU)

	productInfoBySKUChan := make(chan ProductInfo)

	errChan := make(chan error)
	defer close(errChan)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		//рутина пропушивает все айтемы в канал для воркеров
		defer close(skuToGetProductChan)

		for _, item := range items {
			select {
			case <-ctx.Done():
				return
			case skuToGetProductChan <- item.SKU:
			}
		}
	}()

	allProductInfo := make(map[model.SKU]ProductInfo, len(items))

	resultWriterDone := make(chan struct{})
	go func() {
		defer func() {
			resultWriterDone <- struct{}{}
			//были кейсы, когда эта рутина не успевала записать результат в мапу, поэтому добавил этот канал
			//чтобы дождаться ее завершения
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case productInfo, ok := <-productInfoBySKUChan:
				if !ok {
					//проверка на закрытие канала
					//чтобы не записать нулевой productInfo
					return
				}
				allProductInfo[productInfo.sku] = productInfo
			}
		}
	}()

	nw := nWorkers(len(items))
	var errorOccured int32
	wg := &sync.WaitGroup{}
	for i := 0; i < nw; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case sku, ok := <-skuToGetProductChan:
					if !ok {
						return
					}
					productInfo, err := apiCall(ctx, sku)
					if err != nil {
						//добавил проверку на кол-во произошедних ерроров.
						//Логгируем только первую ошибку и отменяем контекст
						//Если писать все ошибки в канал, то будут течь рутины на записи в err канал
						//т.к читающая рутина получает первую ошибку, логирует и завершается.
						if atomic.AddInt32(&errorOccured, 1) > 1 {
							return
						}
						errChan <- err
						return
					}
					select {
					case <-ctx.Done():
						return
					case productInfoBySKUChan <- ProductInfo{
						sku:         sku,
						ProductInfo: productInfo,
					}:
					}
				}
			}
		}()
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-errChan:
				if !ok {
					return
				}
				slog.Error(fmt.Sprintf("%s: error while getting product info for sku: %v", op, err))
				cancel()
				errApiCall = err
				return
			}
		}
	}()

	wg.Wait()
	close(productInfoBySKUChan)
	<-resultWriterDone
	if errApiCall != nil {
		return nil, errApiCall
	}

	return allProductInfo, nil
}
