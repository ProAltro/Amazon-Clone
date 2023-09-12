package entity

type Seller struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

type SellerService interface {
	CreateSeller(seller *Seller) (*Seller, error)
	FindSellerByEmail(email string) (*Seller, error)
	FindSellerByID(id int) (*Seller, error)
	AuthenticateSeller(email string, password string) (*Seller, error)
}
