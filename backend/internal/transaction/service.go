package transaction

import (
	"context"
	"strings"
	"time"
)

const (
	maxAmountCents = 1_000_000_000_000 // 10 billion pesos
	maxNoteLen     = 200
	defaultLimit   = 50
	maxLimit       = 200
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID string, in CreateInput) (Transaction, error) {
	if in.Type != "income" && in.Type != "expense" {
		return Transaction{}, &ValidationError{"Type must be income or expense"}
	}
	if err := validateAmount(in.AmountCents); err != nil {
		return Transaction{}, err
	}
	note, err := cleanNote(in.Note)
	if err != nil {
		return Transaction{}, err
	}
	date, err := parseDate(in.OccurredOn)
	if err != nil {
		return Transaction{}, err
	}

	return s.repo.Create(ctx, userID, NewTransaction{
		AccountID:   in.AccountID,
		CategoryID:  in.CategoryID,
		Type:        in.Type,
		AmountCents: in.AmountCents,
		Note:        note,
		OccurredOn:  date,
	})
}

func (s *Service) List(ctx context.Context, userID string, in ListInput) ([]Transaction, error) {
	if in.Type != "" && in.Type != "income" && in.Type != "expense" {
		return nil, &ValidationError{"Type must be income or expense"}
	}

	f := Filter{Type: in.Type, Limit: in.Limit, Offset: in.Offset}

	if in.From != "" {
		d, err := time.Parse("2006-01-02", in.From)
		if err != nil {
			return nil, &ValidationError{"from must be a date like 2026-11-01"}
		}
		f.From = &d
	}
	if in.To != "" {
		d, err := time.Parse("2006-01-02", in.To)
		if err != nil {
			return nil, &ValidationError{"to must be a date like 2026-11-30"}
		}
		f.To = &d
	}
	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return nil, &ValidationError{"from must be on or before to"}
	}

	if in.AccountID != "" {
		f.AccountID = &in.AccountID
	}
	if in.CategoryID != "" {
		f.CategoryID = &in.CategoryID
	}

	if f.Limit <= 0 {
		f.Limit = defaultLimit
	}
	if f.Limit > maxLimit {
		f.Limit = maxLimit
	}

	return s.repo.List(ctx, userID, f)
}

func (s *Service) Update(ctx context.Context, userID, id string, in UpdateInput) (Transaction, error) {
	if err := validateAmount(in.AmountCents); err != nil {
		return Transaction{}, err
	}
	note, err := cleanNote(in.Note)
	if err != nil {
		return Transaction{}, err
	}
	date, err := parseDate(in.OccurredOn)
	if err != nil {
		return Transaction{}, err
	}

	return s.repo.Update(ctx, userID, id, Changes{
		CategoryID:  in.CategoryID,
		AmountCents: in.AmountCents,
		Note:        note,
		OccurredOn:  date,
	})
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

// ---- helpers ----

func validateAmount(cents int64) error {
	if cents <= 0 {
		return &ValidationError{"Enter an amount greater than 0"}
	}
	if cents > maxAmountCents {
		return &ValidationError{"That amount is too large"}
	}
	return nil
}

// cleanNote trims the note and returns nil when it is empty.
func cleanNote(note string) (*string, error) {
	note = strings.TrimSpace(note)
	if note == "" {
		return nil, nil
	}
	if len([]rune(note)) > maxNoteLen {
		return nil, &ValidationError{"Note must be 200 characters or fewer"}
	}
	return &note, nil
}

// parseDate reads YYYY-MM-DD. An empty value means today.
func parseDate(s string) (time.Time, error) {
	if s == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, &ValidationError{"Date must look like 2026-11-12"}
	}
	return d, nil
}
