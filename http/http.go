package http

import (
	"fmt"
	"time"

	"github.com/ProAltro/Amazon-Clone/mysql"
	"github.com/google/uuid"
)

type HTTPService struct {
	UserService   *mysql.UserService
	SellerService *mysql.SellerService
}

func NewHTTPService(userService *mysql.UserService) *HTTPService {
	return &HTTPService{
		UserService: userService,
	}
}

type session struct {
	username  string
	expiresAt time.Time
}

type sessionsHandler interface {
	// CreateSession creates a new session for an existing user
	CreateSession(email string, expiresAt time.Time) error
	// GetSession returns the enail of the user associated with the given sessionID.
	// If the session is not found or is expired, an error is returned.
	GetSession(sessionID string) (string, error)
	// DeleteSession deletes the session with the given sessionID.
	DeleteSession(sessionID string) error
	// DeleteExpiredSessions deletes all expired sessions from the database.
	DeleteExpiredSessions() error
}

var sessions map[string]session = make(map[string]session)

func CreateSession(email string, expiresAt time.Time) (string, error) {
	//generate random sessionID
	sessionID := uuid.New().String()
	sessions[sessionID] = session{
		username:  email,
		expiresAt: expiresAt,
	}
	for k, v := range sessions {
		fmt.Println(k, v)
	}
	return sessionID, nil
}

func GetSession(sessionID string) (string, error) {
	sessions := sessions
	fmt.Println(sessions)
	s, ok := sessions[sessionID]
	if !ok {
		return "", fmt.Errorf("session not found")
	}
	if s.expiresAt.Before(time.Now()) {
		return "", fmt.Errorf("session expired")
	}
	return s.username, nil
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
			delete(sessions, s.username)
		}
	}
	return nil
}
