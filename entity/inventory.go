package entity

type Inventory struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Seller    int `json:"seller"`
}

type SellerInventory struct {
	Products []Inventory `json:"products"`
}

// amazon clone implementation of inventory
type InventoryService interface {
	CreateInventory(inventory *Inventory) (*Inventory, error)
	FindInventoryByProduct(product int) (*Inventory, error)
	FindInventoryBySeller(seller int) (*SellerInventory, error)
	UpdateInventory(pid, quantity int) (*Inventory, error)
	UnlinkInventory(inventory *Inventory) error
}
