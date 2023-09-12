package entity

import "time"

type User struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	DOJ      time.Time `json:"date_of_join"`
	IsPrime  bool      `json:"is_prime"`
}

type UserService interface {
	CreateUser(user *User) (*User, error)
	FindAllUsers() ([]User, error)
	FindUserByID(id int) (*User, error)
	FindUserByEmail(email string) (*User, error)
	AuthenticateUser(email string, password string) (*User, error)
}
