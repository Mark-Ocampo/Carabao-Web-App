package account

import (
	"context"
	"strings"
)

var validTypes = map[string]bool{
	"cash":    true,
	"bank":    true,
	"ewallet": true,
	"credit":  true,
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID, name, accType string, balanceCents int64) (Account, error) {
	name, err := validateName(name)
	if err != nil {
		return Account{}, err
	}
	if !validTypes[accType] {
		return Account{}, &ValidationError{"Type must be cash, bank, ewallet or credit"}
	}
	return s.repo.Create(ctx, userID, name, accType, balanceCents)
}

func (s *Service) List(ctx context.Context, userID string) ([]Account, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Rename(ctx context.Context, userID, id, name string) (Account, error) {
	name, err := validateName(name)
	if err != nil {
		return Account{}, err
	}
	return s.repo.Rename(ctx, userID, id, name)
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

func validateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", &ValidationError{"Enter an account name"}
	}
	if len([]rune(name)) > 60 {
		return "", &ValidationError{"Account name must be 60 characters or fewer"}
	}
	return name, nil
}
