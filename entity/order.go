package entity

type Order struct {
	Id         int    `json:"id"`
	User       int    `json:"user"`
	TotalPrice int    `json:"total_price"`
	Date       string `json:"date"`
}
