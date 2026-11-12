package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const columns = "id, name, type, balance_cents, created_at"

func scan(row pgx.Row) (Account, error) {
	var a Account
	err := row.Scan(&a.ID, &a.Name, &a.Type, &a.BalanceCents, &a.CreatedAt)
	return a, err
}

func (r *Repository) Create(ctx context.Context, userID, name, accType string, balanceCents int64) (Account, error) {
	return scan(r.db.QueryRow(ctx,
		`INSERT INTO accounts (user_id, name, type, balance_cents)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+columns,
		userID, name, accType, balanceCents,
	))
}

func (r *Repository) List(ctx context.Context, userID string) ([]Account, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+columns+` FROM accounts WHERE user_id = $1 ORDER BY created_at`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		a, err := scan(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *Repository) Rename(ctx context.Context, userID, id, name string) (Account, error) {
	a, err := scan(r.db.QueryRow(ctx,
		`UPDATE accounts SET name = $3 WHERE id = $1 AND user_id = $2 RETURNING `+columns,
		id, userID, name,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return Account{}, ErrNotFound
	}
	return a, err
}

func (r *Repository) Delete(ctx context.Context, userID, id string) error {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT count(*) FROM transactions WHERE account_id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrInUse
	}

	tag, err := r.db.Exec(ctx,
		`DELETE FROM accounts WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
