package account

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Mark-Ocampo/carabao-web-app/api/internal/httpx"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/server/middleware"
)

type Handler struct {
	svc   *Service
	authn middleware.Authenticator
}

func NewHandler(svc *Service, authn middleware.Authenticator) *Handler {
	return &Handler{svc: svc, authn: authn}
}

func (h *Handler) Routes(r chi.Router) {
	r.Use(middleware.RequireAuth(h.authn))
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Patch("/{id}", h.rename)
	r.Delete("/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.List(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, accounts)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		BalanceCents int64  `json:"balanceCents"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	a, err := h.svc.Create(r.Context(), middleware.UserID(r.Context()), in.Name, in.Type, in.BalanceCents)
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, a)
}

func (h *Handler) rename(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Account not found")
		return
	}

	var in struct {
		Name string `json:"name"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	a, err := h.svc.Rename(r.Context(), middleware.UserID(r.Context()), id, in.Name)
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, a)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Account not found")
		return
	}

	if err := h.svc.Delete(r.Context(), middleware.UserID(r.Context()), id); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		httpx.Error(w, http.StatusBadRequest, ve.Message)
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "Account not found")
	case errors.Is(err, ErrInUse):
		httpx.Error(w, http.StatusConflict, "This account has transactions and can't be deleted")
	default:
		log.Printf("account: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong")
	}
}
