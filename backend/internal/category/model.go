package category

import "errors"

// Kind is the wants vs. needs tag and only applies to expense categories.
type Category struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"`
	Kind *string `json:"kind"`
}

var (
	ErrNotFound  = errors.New("category not found")
	ErrDuplicate = errors.New("category already exists")
)

// ValidationError is a problem with the input that the user can fix.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
