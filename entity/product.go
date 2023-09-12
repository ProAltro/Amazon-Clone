package entity

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Features    string `json:"features"`
	Seller      int    `json:"seller"`
}
