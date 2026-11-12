package category

import (
	"context"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID, name, catType string, kind *string) (Category, error) {
	name, err := validateName(name)
	if err != nil {
		return Category{}, err
	}
	if catType != "income" && catType != "expense" {
		return Category{}, &ValidationError{"Type must be income or expense"}
	}
	kind, err = validateKind(catType, kind)
	if err != nil {
		return Category{}, err
	}
	return s.repo.Create(ctx, userID, name, catType, kind)
}

func (s *Service) List(ctx context.Context, userID, filterType string) ([]Category, error) {
	if filterType != "" && filterType != "income" && filterType != "expense" {
		return nil, &ValidationError{"Type must be income or expense"}
	}
	return s.repo.List(ctx, userID, filterType)
}

func (s *Service) Update(ctx context.Context, userID, id, name string, kind *string) (Category, error) {
	existing, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return Category{}, err
	}
	name, err = validateName(name)
	if err != nil {
		return Category{}, err
	}
	kind, err = validateKind(existing.Type, kind)
	if err != nil {
		return Category{}, err
	}
	return s.repo.Update(ctx, userID, id, name, kind)
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

func validateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", &ValidationError{"Enter a category name"}
	}
	if len([]rune(name)) > 60 {
		return "", &ValidationError{"Category name must be 60 characters or fewer"}
	}
	return name, nil
}

// Expense categories must be tagged need or want. Income categories have no kind.
func validateKind(catType string, kind *string) (*string, error) {
	if catType == "income" {
		return nil, nil
	}
	if kind == nil || (*kind != "need" && *kind != "want") {
		return nil, &ValidationError{"Choose need or want for expense categories"}
	}
	return kind, nil
}
