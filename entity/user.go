package entity

import "context"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	DOJ      string `json:"date_of_join"`
}

type UserService interface {
	CreateUser(ctx context.Context, name string, email string, password string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	AuthenticateUser(ctx context.Context, string, password string) (*User, error)
}
