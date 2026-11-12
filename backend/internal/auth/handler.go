package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Mark-Ocampo/carabao-web-app/api/internal/httpx"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/server/middleware"
)

type Handler struct {
	svc          *Service
	secureCookie bool
}

func NewHandler(svc *Service, secureCookie bool) *Handler {
	return &Handler{svc: svc, secureCookie: secureCookie}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/signup", h.signup)
	r.Post("/login", h.login)
	r.Post("/logout", h.logout)
	r.With(middleware.RequireAuth(h.svc)).Get("/me", h.me)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, session, err := h.svc.Signup(r.Context(), in.Email, in.Password, in.FullName)
	if err != nil {
		h.fail(w, err)
		return
	}

	h.setCookie(w, session)
	httpx.JSON(w, http.StatusCreated, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, session, err := h.svc.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		h.fail(w, err)
		return
	}

	h.setCookie(w, session)
	httpx.JSON(w, http.StatusOK, user)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(middleware.SessionCookie); err == nil {
		if err := h.svc.Logout(r.Context(), c.Value); err != nil {
			log.Printf("auth: logout: %v", err)
		}
	}
	h.clearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, err := h.svc.Me(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, user)
}

// ---- helpers ----

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		httpx.Error(w, http.StatusBadRequest, ve.Message)
	case errors.Is(err, ErrEmailTaken):
		httpx.Error(w, http.StatusConflict, "This email is already registered")
	case errors.Is(err, ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, "Invalid email or password")
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusUnauthorized, "Not signed in")
	default:
		log.Printf("auth: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong")
	}
}

func (h *Handler) setCookie(w http.ResponseWriter, s Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    s.Token,
		Path:     "/",
		Expires:  s.ExpiresAt,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}
