package entity

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Features    string `json:"features" binding:"required"`
	Seller      int    `json:"seller" binding:"required"`
}

// amazon clone implementation of products
type ProductService interface {
	CreateProduct(product *Product) (*Product, error)
	FindAllProducts() ([]Product, error)
	FindProductByID(id int) (*Product, error)
	FindProductByName(name string) (*Product, error)
	FindProductsBySeller(seller int) ([]Product, error)
	FindProductsByFilter(filter *ProductFilter) ([]Product, error)
	UpdateProduct(product *Product) (*Product, error)
}

// filter struct
type ProductFilter struct {
	MinPrice int   `json:"min_price"`
	MaxPrice int   `json:"max_price"`
	Sellers  []int `json:"sellers"`
}

// update struct
