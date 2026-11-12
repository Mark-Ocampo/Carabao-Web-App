package category

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
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.List(r.Context(), middleware.UserID(r.Context()), r.URL.Query().Get("type"))
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, categories)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string  `json:"name"`
		Type string  `json:"type"`
		Kind *string `json:"kind"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	c, err := h.svc.Create(r.Context(), middleware.UserID(r.Context()), in.Name, in.Type, in.Kind)
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, c)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Category not found")
		return
	}

	var in struct {
		Name string  `json:"name"`
		Kind *string `json:"kind"`
	}
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	c, err := h.svc.Update(r.Context(), middleware.UserID(r.Context()), id, in.Name, in.Kind)
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, c)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Category not found")
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
		httpx.Error(w, http.StatusNotFound, "Category not found")
	case errors.Is(err, ErrDuplicate):
		httpx.Error(w, http.StatusConflict, "You already have a category with that name")
	default:
		log.Printf("category: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong")
	}
}
