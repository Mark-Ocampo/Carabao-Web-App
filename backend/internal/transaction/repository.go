package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const columns = "id, account_id, category_id, type, amount_cents, note, occurred_on, created_at"

func scan(row pgx.Row) (Transaction, error) {
	var t Transaction
	var date time.Time
	err := row.Scan(&t.ID, &t.AccountID, &t.CategoryID, &t.Type, &t.AmountCents, &t.Note, &date, &t.CreatedAt)
	t.OccurredOn = date.Format("2006-01-02")
	return t, err
}

// signed is how much a transaction moves the account balance.
func signed(txType string, cents int64) int64 {
	if txType == "expense" {
		return -cents
	}
	return cents
}

// checkCategory makes sure the category is the user's and matches the transaction type.
func checkCategory(ctx context.Context, tx pgx.Tx, userID string, categoryID *string, txType string) error {
	if categoryID == nil {
		return nil
	}
	var catType string
	err := tx.QueryRow(ctx,
		`SELECT type FROM categories WHERE id = $1 AND user_id = $2`,
		*categoryID, userID,
	).Scan(&catType)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidCategory
	}
	if err != nil {
		return err
	}
	if catType != txType {
		return ErrCategoryMismatch
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, userID string, n NewTransaction) (Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Transaction{}, err
	}
	defer tx.Rollback(ctx)

	// Lock the account row so concurrent requests can't corrupt the balance
	var one int
	err = tx.QueryRow(ctx,
		`SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2 FOR UPDATE`,
		n.AccountID, userID,
	).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return Transaction{}, ErrInvalidAccount
	}
	if err != nil {
		return Transaction{}, err
	}

	if err := checkCategory(ctx, tx, userID, n.CategoryID, n.Type); err != nil {
		return Transaction{}, err
	}

	t, err := scan(tx.QueryRow(ctx,
		`INSERT INTO transactions (user_id, account_id, category_id, type, amount_cents, note, occurred_on)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING `+columns,
		userID, n.AccountID, n.CategoryID, n.Type, n.AmountCents, n.Note, n.OccurredOn,
	))
	if err != nil {
		return Transaction{}, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance_cents = balance_cents + $2 WHERE id = $1`,
		n.AccountID, signed(n.Type, n.AmountCents),
	)
	if err != nil {
		return Transaction{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Transaction{}, err
	}
	return t, nil
}

func (r *Repository) List(ctx context.Context, userID string, f Filter) ([]Transaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+columns+` FROM transactions
		 WHERE user_id = $1
		   AND ($2::date IS NULL OR occurred_on >= $2::date)
		   AND ($3::date IS NULL OR occurred_on <= $3::date)
		   AND ($4::uuid IS NULL OR account_id = $4::uuid)
		   AND ($5::uuid IS NULL OR category_id = $5::uuid)
		   AND ($6::text = '' OR type = $6::text)
		 ORDER BY occurred_on DESC, created_at DESC
		 LIMIT $7 OFFSET $8`,
		userID, f.From, f.To, f.AccountID, f.CategoryID, f.Type, f.Limit, f.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		t, err := scan(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

func (r *Repository) Update(ctx context.Context, userID, id string, c Changes) (Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Transaction{}, err
	}
	defer tx.Rollback(ctx)

	var accountID, txType string
	var oldAmount int64
	err = tx.QueryRow(ctx,
		`SELECT account_id, type, amount_cents FROM transactions
		 WHERE id = $1 AND user_id = $2 FOR UPDATE`,
		id, userID,
	).Scan(&accountID, &txType, &oldAmount)
	if errors.Is(err, pgx.ErrNoRows) {
		return Transaction{}, ErrNotFound
	}
	if err != nil {
		return Transaction{}, err
	}

	if err := checkCategory(ctx, tx, userID, c.CategoryID, txType); err != nil {
		return Transaction{}, err
	}

	t, err := scan(tx.QueryRow(ctx,
		`UPDATE transactions
		 SET category_id = $3, amount_cents = $4, note = $5, occurred_on = $6
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+columns,
		id, userID, c.CategoryID, c.AmountCents, c.Note, c.OccurredOn,
	))
	if err != nil {
		return Transaction{}, err
	}

	delta := signed(txType, c.AmountCents) - signed(txType, oldAmount)
	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance_cents = balance_cents + $2 WHERE id = $1`,
		accountID, delta,
	)
	if err != nil {
		return Transaction{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Transaction{}, err
	}
	return t, nil
}

func (r *Repository) Delete(ctx context.Context, userID, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var accountID, txType string
	var amount int64
	err = tx.QueryRow(ctx,
		`SELECT account_id, type, amount_cents FROM transactions
		 WHERE id = $1 AND user_id = $2 FOR UPDATE`,
		id, userID,
	).Scan(&accountID, &txType, &amount)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id); err != nil {
		return err
	}

	// Reverse the effect on the balance
	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance_cents = balance_cents - $2 WHERE id = $1`,
		accountID, signed(txType, amount),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
