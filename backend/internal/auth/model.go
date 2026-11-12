package auth

import (
	"errors"
	"time"
)

type User struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	FullName           string    `json:"fullName"`
	EmergencyMinMonths int       `json:"emergencyMinMonths"`
	CreatedAt          time.Time `json:"createdAt"`
}

type Session struct {
	Token     string
	ExpiresAt time.Time
}

var (
	ErrNotFound           = errors.New("not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
)

// ValidationError is a problem with the input that the user can fix.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
