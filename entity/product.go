package entity

import "context"

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       int      `json:"price" binding:"required"`
	Seller      string   `json:"seller" binding:"required"`
	Images      []string `json:"images"`
}

type ProductService interface {
	CreateProduct(ctx context.Context, name string, description string, price int, seller string, images []string) (*Product, error)
	GetProduct(ctx context.Context, id int) (*Product, error)
	GetProducts(ctx context.Context, ids []int) ([]Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	DeleteProduct(ctx context.Context, id int) error
}

func (p *Product) GetImageURLS() []string {
	image_urls := []string{}
	for _, image := range p.Images {
		image_url := "https://pramitpal.me/amazon/api/v1/images/" + image + ".jpg"
		image_urls = append(image_urls, image_url)
	}
	return image_urls
}

func (p *Product) Copy() Product {
	images := []string{}
	for _, image := range p.Images {
		images = append(images, image)
	}

	return Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Seller:      p.Seller,
		Images:      images,
	}
}

func (p *Product) Update(new Product) (Product, error) {
	if new.Name != "" {
		p.Name = new.Name
	}
	if new.Description != "" {
		p.Description = new.Description
	}
	if new.Price != 0 {
		p.Price = new.Price
	}
	if new.Seller != "" {
		p.Seller = new.Seller
	}
	if new.Images != nil {
		p.Images = new.Images
	}
	return *p, nil
}
