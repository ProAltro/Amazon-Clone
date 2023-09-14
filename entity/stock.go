package entity

import "context"

type Stock struct {
	Product  Product `json:"product_id"`
	Quantity int     `json:"quantity"`
}

type InventoryService interface {
	AddStockToInventory(ctx context.Context, id int, quantity int) error
	RemoveStockFromInventory(ctx context.Context, id int) error
	GetStockFromInventory(ctx context.Context, id int) (*Stock, error)
	GetAllStocksFromInventory(ctx context.Context) ([]*Stock, error)
	UpdateStockInInventory(ctx context.Context, id int, quantity int) error
}
