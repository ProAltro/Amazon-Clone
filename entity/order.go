package entity

import (
	"context"
)

type Order struct {
	ID       int     `json:"id"`
	Products []Stock `json:"products" binding:"required"`
	Total    int     `json:"total"`
	UID      int     `json:"user" binding:"required"`
}

type OrderService interface {
	CreateOrder(ctx context.Context, products []Stock, total int) (*Order, error)
	GetOrder(ctx context.Context, id int) (*Order, error)
	GetOrders(ctx context.Context, ids []int) ([]Order, error)
	GetOrdersOfUser(ctx context.Context, uid int) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	DeleteOrder(ctx context.Context, id int) error
}
