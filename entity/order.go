package entity

type Order struct {
	Id         int    `json:"id"`
	User       int    `json:"user"`
	TotalPrice int    `json:"total_price"`
	Date       string `json:"date"`
}

type OrderService interface {
	CreateOrder(order *Order) (*Order, error)
	FindAllOrders() ([]Order, error)
	FindOrderById(id int) (*Order, error)
	FindOrdersByUser(user int) ([]Order, error)
}
