package auth

import (
	"context"
	"errors"
	"time"

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

const userColumns = "id, email, full_name, emergency_min_months, created_at"

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, fullName string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, full_name)
		 VALUES ($1, $2, $3)
		 RETURNING `+userColumns,
		email, passwordHash, fullName,
	).Scan(&u.ID, &u.Email, &u.FullName, &u.EmergencyMinMonths, &u.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return User{}, ErrEmailTaken
	}
	return u, err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (User, string, error) {
	var u User
	var hash string
	err := r.db.QueryRow(ctx,
		`SELECT `+userColumns+`, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.FullName, &u.EmergencyMinMonths, &u.CreatedAt, &hash)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrNotFound
	}
	return u, hash, err
}

func (r *Repository) FindByID(ctx context.Context, id string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.FullName, &u.EmergencyMinMonths, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (r *Repository) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

// FindSessionUser returns the user ID for a valid, unexpired session.
func (r *Repository) FindSessionUser(ctx context.Context, tokenHash string) (string, error) {
	var userID string
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > now()`,
		tokenHash,
	).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return userID, err
}

func (r *Repository) DeleteSession(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}
