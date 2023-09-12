package entity

import (
	"database/sql/driver"
	"time"
)

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

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	// if value is [8]uint convert it to string
	var val string
	switch value.(type) {
	case []uint8:
		val = string(value.([]uint8))
	default:
		val = ""
	}
	datetime, err := time.Parse("2006-01-02 15:04:05", val)
	nt.Time = datetime
	nt.Valid = err == nil
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
