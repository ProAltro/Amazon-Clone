package http

import (
	"fmt"
	"time"

	"github.com/ProAltro/Amazon-Clone/mysql"
	"github.com/google/uuid"
)

type HTTPService struct {
	UserService      *mysql.UserService
	ProductService   *mysql.ProductService
	CartService      *mysql.CartService
	OrderService     *mysql.OrderService
	InventoryService *mysql.InventoryService
}

func NewHTTPService(userService *mysql.UserService, productService *mysql.ProductService, cartService *mysql.CartService, orderService *mysql.OrderService, inventoryService *mysql.InventoryService) *HTTPService {
	return &HTTPService{
		UserService:      userService,
		ProductService:   productService,
		CartService:      cartService,
		OrderService:     orderService,
		InventoryService: inventoryService,
	}
}

type session struct {
	uid       int
	email     string
	expiresAt time.Time
	sessionID string
}

var sessions map[string]session = make(map[string]session)

func CreateSession(uid int, email string, expiresAt time.Time) (string, error) {
	//generate random sessionID
	sessionID := uuid.New().String()
	sessions[sessionID] = session{
		uid:       uid,
		email:     email,
		expiresAt: expiresAt,
		sessionID: sessionID,
	}

	return sessionID, nil
}

func GetSession(sessionID string) (string, int, error) {
	sessions := sessions
	s, ok := sessions[sessionID]
	if !ok {
		return "", -1, fmt.Errorf("session not found")
	}
	if s.expiresAt.Before(time.Now()) {
		return "", -1, fmt.Errorf("session expired")
	}
	return s.email, s.uid, nil
}

func DeleteSession(sessionID string) error {
	sessions := sessions
	delete(sessions, sessionID)
	return nil
}

func DeleteExpiredSessions() error {
	sessions := sessions
	for _, s := range sessions {
		if s.expiresAt.Before(time.Now()) {
			delete(sessions, s.sessionID)
		}
	}
	return nil
}
