package http

import "github.com/ProAltro/Amazon-Clone/mysql"

type HTTPService struct {
	UserService *mysql.UserService
}

func NewHTTPService(userService *mysql.UserService) *HTTPService {
	return &HTTPService{
		UserService: userService,
	}
}
