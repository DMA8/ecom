package cart

import (
	"ecom/cart/internal/model"
)

func (c *Cart) IsSKUValid(itemInfo model.ProductInfo) error {
	if itemInfo.Name == "" || itemInfo.Price == 0 {
		return model.ErrProductServiceSuspiciousResponse
	}
	return nil
}
