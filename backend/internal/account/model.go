package account

import (
	"errors"
	"time"
)

type Account struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	BalanceCents int64     `json:"balanceCents"`
	CreatedAt    time.Time `json:"createdAt"`
}

var (
	ErrNotFound = errors.New("account not found")
	ErrInUse    = errors.New("account has transactions")
)

// ValidationError is a problem with the input that the user can fix.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
