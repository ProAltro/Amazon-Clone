package entity

import (
	"context"
	"encoding/json"
)

type Stock struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type InventoryService interface {
	AddStockToInventory(ctx context.Context, id int, quantity int) error
	RemoveStockFromInventory(ctx context.Context, id int) error
	GetStockFromInventory(ctx context.Context, id int) (*Stock, error)
	GetAllStocksFromInventory(ctx context.Context) ([]*Stock, error)
	UpdateStockInInventory(ctx context.Context, id int, quantity int) error
}

func ToJSON(stocks []Stock) ([]byte, error) {
	var products []map[string]interface{}
	for _, stock := range stocks {
		product := map[string]interface{}{
			"id":          stock.Product.ID,
			"name":        stock.Product.Name,
			"description": stock.Product.Description,
			"price":       stock.Product.Price,
			"seller":      stock.Product.Seller,
		}
		stockMap := map[string]interface{}{
			"product":  product,
			"quantity": stock.Quantity,
		}
		products = append(products, stockMap)
	}
	return json.Marshal(products)
}

func FromJSON(data []byte) ([]Stock, error) {
	var stocks []Stock
	var products []map[string]interface{}
	err := json.Unmarshal(data, &products)
	if err != nil {
		return nil, err
	}
	for _, stock := range products {
		product := stock["product"].(map[string]interface{})
		quantity := stock["quantity"].(float64)
		stocks = append(stocks, Stock{
			Product: Product{
				ID:          int(product["id"].(float64)),
				Name:        product["name"].(string),
				Description: product["description"].(string),
				Price:       int(product["price"].(float64)),
				Seller:      product["seller"].(string),
			},
			Quantity: int(quantity),
		})
	}
	return stocks, nil
}
