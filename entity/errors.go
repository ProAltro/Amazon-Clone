package entity

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal error")
var ErrServer = fmt.Errorf("server error: %w", ErrInternal)
var ErrDB = fmt.Errorf("database error: %w", ErrServer)
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrForbidden = errors.New("forbidden")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInvalidInput = errors.New("invalid input")
var ErrBadData = errors.New("bad data")

func GetStatusCode(err error) int {
	if err == nil {
		return 200
	} else if errors.Is(err, ErrInternal) {
		return 500
	} else if errors.Is(err, ErrServer) {
		return 500
	} else if errors.Is(err, ErrDB) {
		return 500
	} else if errors.Is(err, ErrNotFound) {
		return 404
	} else if errors.Is(err, ErrConflict) {
		return 409
	} else if errors.Is(err, ErrForbidden) {
		return 403
	} else if errors.Is(err, ErrUnauthorized) {
		return 401
	} else if errors.Is(err, ErrInvalidInput) {
		return 400
	} else if errors.Is(err, ErrBadData) {
		return 400
	} else {
		return 500
	}
}
