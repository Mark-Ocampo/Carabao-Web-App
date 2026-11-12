package transaction

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	CategoryID  *string   `json:"categoryId"`
	Type        string    `json:"type"`
	AmountCents int64     `json:"amountCents"`
	Note        *string   `json:"note"`
	OccurredOn  string    `json:"occurredOn"` // YYYY-MM-DD
	CreatedAt   time.Time `json:"createdAt"`
}

// ---- input from the handler (raw strings) ----

type CreateInput struct {
	AccountID   string
	CategoryID  *string
	Type        string
	AmountCents int64
	Note        string
	OccurredOn  string
}

type UpdateInput struct {
	CategoryID  *string
	AmountCents int64
	Note        string
	OccurredOn  string
}

type ListInput struct {
	From       string
	To         string
	AccountID  string
	CategoryID string
	Type       string
	Limit      int
	Offset     int
}

// ---- validated values passed to the repository ----

type NewTransaction struct {
	AccountID   string
	CategoryID  *string
	Type        string
	AmountCents int64
	Note        *string
	OccurredOn  time.Time
}

type Changes struct {
	CategoryID  *string
	AmountCents int64
	Note        *string
	OccurredOn  time.Time
}

type Filter struct {
	From       *time.Time
	To         *time.Time
	AccountID  *string
	CategoryID *string
	Type       string
	Limit      int
	Offset     int
}

var (
	ErrNotFound         = errors.New("transaction not found")
	ErrInvalidAccount   = errors.New("invalid account")
	ErrInvalidCategory  = errors.New("invalid category")
	ErrCategoryMismatch = errors.New("category type does not match transaction type")
)

// ValidationError is a problem with the input that the user can fix.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
