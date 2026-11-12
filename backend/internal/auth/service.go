package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	sessionTTL = 30 * 24 * time.Hour

	// argon2id parameters (OWASP recommended minimum)
	argonMemory  = 19 * 1024 // KiB
	argonTime    = 2
	argonThreads = 1
	argonKeyLen  = 32
	argonSaltLen = 16
)

// dummyHash is compared against when an email is unknown,
// so login takes about the same time either way.
var dummyHash, _ = hashPassword("not-a-real-password")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Signup(ctx context.Context, email, password, fullName string) (User, Session, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	fullName = strings.TrimSpace(fullName)

	if err := validateSignup(email, password, fullName); err != nil {
		return User{}, Session{}, err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return User{}, Session{}, err
	}

	user, err := s.repo.CreateUser(ctx, email, hash, fullName)
	if err != nil {
		return User{}, Session{}, err
	}

	session, err := s.newSession(ctx, user.ID)
	if err != nil {
		return User{}, Session{}, err
	}
	return user, session, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (User, Session, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, hash, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		verifyPassword(dummyHash, password)
		return User{}, Session{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, Session{}, err
	}

	if !verifyPassword(hash, password) {
		return User{}, Session{}, ErrInvalidCredentials
	}

	session, err := s.newSession(ctx, user.ID)
	if err != nil {
		return User{}, Session{}, err
	}
	return user, session, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, hashToken(token))
}

// Authenticate returns the user ID for a session token.
// It satisfies middleware.Authenticator.
func (s *Service) Authenticate(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", ErrUnauthorized
	}
	userID, err := s.repo.FindSessionUser(ctx, hashToken(token))
	if errors.Is(err, ErrNotFound) {
		return "", ErrUnauthorized
	}
	return userID, err
}

func (s *Service) Me(ctx context.Context, userID string) (User, error) {
	return s.repo.FindByID(ctx, userID)
}

// ---- helpers ----

func validateSignup(email, password, fullName string) error {
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return &ValidationError{"Enter a valid email address"}
	}
	if fullName == "" {
		return &ValidationError{"Enter your name"}
	}
	n := len([]rune(password))
	if n < 8 {
		return &ValidationError{"Password must be at least 8 characters"}
	}
	if n > 128 {
		return &ValidationError{"Password must be 128 characters or fewer"}
	}
	return nil
}

func (s *Service) newSession(ctx context.Context, userID string) (Session, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return Session{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	expires := time.Now().Add(sessionTTL)

	// Only a hash of the token is stored, so a database leak can't be used to sign in
	if err := s.repo.CreateSession(ctx, userID, hashToken(token), expires); err != nil {
		return Session{}, err
	}
	return Session{Token: token, ExpiresAt: expires}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func verifyPassword(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}

	var memory, iterations uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &threads); err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	got := argon2.IDKey([]byte(password), salt, iterations, memory, threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1
}
