package middleware

import (
	"context"
	"net/http"

	"github.com/Mark-Ocampo/carabao-web-app/api/internal/httpx"
)

// SessionCookie is the name of the cookie that carries the session token.
const SessionCookie = "cara_session"

type ctxKey struct{}

// Authenticator turns a session token into a user ID.
// The auth service implements this.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (string, error)
}

// RequireAuth rejects requests without a valid session and puts the user ID in the context.
func RequireAuth(a Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(SessionCookie)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "not signed in")
				return
			}

			userID, err := a.Authenticate(r.Context(), c.Value)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "not signed in")
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserID returns the signed-in user's ID, or "" if there is none.
func UserID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}
