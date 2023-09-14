package entity

import "context"

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Seller      string `json:"seller" binding:"required"`
}

type ProductService interface {
	CreateProduct(ctx context.Context, name string, description string, price int, seller string) (*Product, error)
	GetProduct(ctx context.Context, id int) (*Product, error)
	GetProducts(ctx context.Context, ids []int) ([]Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
