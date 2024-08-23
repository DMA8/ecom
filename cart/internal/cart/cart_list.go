package cart

import (
	"context"
	"ecom/cart/internal/model"
	"fmt"
	"log/slog"

	errgroup "golang.org/x/sync/errgroup"
)

func nWorkers(n int) int {
	if n > 3 {
		return 3
	}
	return n
}

func (c *Cart) CartList(ctx context.Context, user model.User) (model.UserCart, error) {
	var totalPrice uint32
	var userCart model.UserCart
	const op = `cart.CartList`

	allItems, err := c.repo.ItemsByUser(ctx, user)
	if err != nil {
		return model.UserCart{}, fmt.Errorf("%s err: %w", op, err)
	}

	if len(allItems) == 0 {
		return model.UserCart{}, nil
	}

	userCart.Items = make([]model.ItemInCart, 0, len(allItems))

	allProductInfo, err := getProductForItemsAsync(ctx, allItems, c.productServiceAPIClient.GetProduct)
	if err != nil {
		return model.UserCart{}, fmt.Errorf("%s err: %w", op, err)
	}

	for _, item := range allItems {
		productInfo := allProductInfo[item.SKU]
		totalPrice += productInfo.Price * uint32(item.Count)
		userCart.Items = append(userCart.Items, model.ItemInCart{
			Sku:   item.SKU,
			Name:  productInfo.Name,
			Price: productInfo.Price,
			Count: item.Count,
		})
	}
	userCart.TotalPrice = totalPrice
	return userCart, nil
}

type ProductInfo struct {
	sku model.SKU
	model.ProductInfo
}

func getProductForItemsErrGroups(ctx context.Context, items []model.Item, apiCall func(ctx context.Context, sku model.SKU) (model.ProductInfo, error)) (map[model.SKU]ProductInfo, error) {
	if len(items) < 1 {
		return make(map[model.SKU]ProductInfo), nil
	}

	//Решил попробовать errgroup
	//По задумке - все рутины должны завершиться, если контекст отменится, или функция ерроргруппы вернет ошибку
	g, ctx := errgroup.WithContext(ctx)

	skuToGetProductChan := make(chan model.SKU)
	productInfoBySKUChan := make(chan ProductInfo)
	defer close(productInfoBySKUChan)

	go func() {
		//пишем все айтемы в синхронный канал, т.к смысла в буф не вижу
		//если бы мы доставали из базы по 1 айтему, то стоило бы делать буф канал
		defer close(skuToGetProductChan)
		for _, item := range items {
			//сначала писал в канал без селекта, но ликдетектор рутин (в cart_list_test) подсказал
			select {
			case <-ctx.Done():
				return
			case skuToGetProductChan <- item.SKU:
			}
		}
	}()

	allProductInfo := make(map[model.SKU]ProductInfo, len(items))

	go func() {
		// в этой рутинке аккумулируем все результаты работы воркеров, которые ходят в productService
		// канал productInfoBySKUChan будет всегда закрыт через дефер, так что эта рутинка не потечет
		for productInfo := range productInfoBySKUChan {
			allProductInfo[productInfo.sku] = productInfo
		}
	}()

	//Не уверен, что оптимально подобрал кол-во воркеров (максимум 3)
	//Думаю, для определения оптимального кол-ва воркеров, нужно смотреть на среднее число айтемов в корзине
	//А еще было бы лучше сразу всю корзину отправлять в productService (думаю для учебных целей нам порезали ходить в productService батчами)
	nw := nWorkers(len(items))

	for i := 0; i < nw; i++ {
		//раним воркеры, которые будут делать запросы в productService
		//решил принимать метод вызова API через аргумент функции, чтобы подсовывать фейки в тестах
		g.Go(func() error {
			for sku := range skuToGetProductChan {
				productInfo, err := apiCall(ctx, sku)
				if err != nil {
					slog.Error(fmt.Sprintf("error while getting product info for sku %d: %v", sku, err))
					return fmt.Errorf("error: %w", err)
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case productInfoBySKUChan <- ProductInfo{sku: sku, ProductInfo: productInfo}:
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return allProductInfo, nil
}
