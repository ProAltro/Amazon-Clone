package entity

import "context"

type Cart struct {
	Products []Stock `json:"products" binding:"required"`
	Total    int     `json:"total" binding:"required"`
	UID      int     `json:"user" binding:"required"`
}

type CartService interface {
	AddProductToCart(ctx context.Context, int, quantity int) error
	RemoveProductFromCart(ctx context.Context, id int) error
	GetCart(ctx context.Context) (*Cart, error)
	ModifyCart(ctx context.Context, id int, quantity int) (*Cart, error)
	Total(ctx context.Context) (int, error)
	ClearCart(ctx context.Context) error
	Checkout(ctx context.Context) error
}
