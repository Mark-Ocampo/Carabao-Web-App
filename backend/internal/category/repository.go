package category

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const columns = "id, name, type, kind"

func scan(row pgx.Row) (Category, error) {
	var c Category
	err := row.Scan(&c.ID, &c.Name, &c.Type, &c.Kind)
	return c, err
}

func mapErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return ErrDuplicate
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *Repository) Create(ctx context.Context, userID, name, catType string, kind *string) (Category, error) {
	c, err := scan(r.db.QueryRow(ctx,
		`INSERT INTO categories (user_id, name, type, kind)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+columns,
		userID, name, catType, kind,
	))
	return c, mapErr(err)
}

// List returns the user's categories. An empty filterType returns all of them.
func (r *Repository) List(ctx context.Context, userID, filterType string) ([]Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+columns+` FROM categories
		 WHERE user_id = $1 AND ($2::text = '' OR type = $2)
		 ORDER BY type, name`,
		userID, filterType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		c, err := scan(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *Repository) Get(ctx context.Context, userID, id string) (Category, error) {
	c, err := scan(r.db.QueryRow(ctx,
		`SELECT `+columns+` FROM categories WHERE id = $1 AND user_id = $2`,
		id, userID,
	))
	return c, mapErr(err)
}

func (r *Repository) Update(ctx context.Context, userID, id, name string, kind *string) (Category, error) {
	c, err := scan(r.db.QueryRow(ctx,
		`UPDATE categories SET name = $3, kind = $4
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+columns,
		id, userID, name, kind,
	))
	return c, mapErr(err)
}

func (r *Repository) Delete(ctx context.Context, userID, id string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM categories WHERE id = $1 AND user_id = $2`,
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
