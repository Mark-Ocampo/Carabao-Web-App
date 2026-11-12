package transaction

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	in := ListInput{
		From:       q.Get("from"),
		To:         q.Get("to"),
		AccountID:  q.Get("accountId"),
		CategoryID: q.Get("categoryId"),
		Type:       q.Get("type"),
	}

	if in.AccountID != "" && !httpx.IsUUID(in.AccountID) {
		httpx.Error(w, http.StatusBadRequest, "accountId is not valid")
		return
	}
	if in.CategoryID != "" && !httpx.IsUUID(in.CategoryID) {
		httpx.Error(w, http.StatusBadRequest, "categoryId is not valid")
		return
	}

	var ok bool
	if in.Limit, ok = intParam(q.Get("limit")); !ok {
		httpx.Error(w, http.StatusBadRequest, "limit must be a positive number")
		return
	}
	if in.Offset, ok = intParam(q.Get("offset")); !ok {
		httpx.Error(w, http.StatusBadRequest, "offset must be a positive number")
		return
	}

	transactions, err := h.svc.List(r.Context(), middleware.UserID(r.Context()), in)
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, transactions)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccountID   string  `json:"accountId"`
		CategoryID  *string `json:"categoryId"`
		Type        string  `json:"type"`
		AmountCents int64   `json:"amountCents"`
		Note        string  `json:"note"`
		OccurredOn  string  `json:"occurredOn"`
	}
	if err := httpx.Decode(r, &body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !httpx.IsUUID(body.AccountID) {
		httpx.Error(w, http.StatusBadRequest, "Choose a valid account")
		return
	}
	if body.CategoryID != nil && !httpx.IsUUID(*body.CategoryID) {
		httpx.Error(w, http.StatusBadRequest, "Choose a valid category")
		return
	}

	t, err := h.svc.Create(r.Context(), middleware.UserID(r.Context()), CreateInput{
		AccountID:   body.AccountID,
		CategoryID:  body.CategoryID,
		Type:        body.Type,
		AmountCents: body.AmountCents,
		Note:        body.Note,
		OccurredOn:  body.OccurredOn,
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, t)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Transaction not found")
		return
	}

	var body struct {
		CategoryID  *string `json:"categoryId"`
		AmountCents int64   `json:"amountCents"`
		Note        string  `json:"note"`
		OccurredOn  string  `json:"occurredOn"`
	}
	if err := httpx.Decode(r, &body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.CategoryID != nil && !httpx.IsUUID(*body.CategoryID) {
		httpx.Error(w, http.StatusBadRequest, "Choose a valid category")
		return
	}

	t, err := h.svc.Update(r.Context(), middleware.UserID(r.Context()), id, UpdateInput{
		CategoryID:  body.CategoryID,
		AmountCents: body.AmountCents,
		Note:        body.Note,
		OccurredOn:  body.OccurredOn,
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !httpx.IsUUID(id) {
		httpx.Error(w, http.StatusNotFound, "Transaction not found")
		return
	}

	if err := h.svc.Delete(r.Context(), middleware.UserID(r.Context()), id); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- helpers ----

// intParam reads an optional non-negative integer query value.
func intParam(s string) (int, bool) {
	if s == "" {
		return 0, true
	}
	n, err := strconv.Atoi(s)
	return n, err == nil && n >= 0
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		httpx.Error(w, http.StatusBadRequest, ve.Message)
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "Transaction not found")
	case errors.Is(err, ErrInvalidAccount):
		httpx.Error(w, http.StatusBadRequest, "Choose a valid account")
	case errors.Is(err, ErrInvalidCategory):
		httpx.Error(w, http.StatusBadRequest, "Choose a valid category")
	case errors.Is(err, ErrCategoryMismatch):
		httpx.Error(w, http.StatusBadRequest, "The category type must match the transaction type")
	default:
		log.Printf("transaction: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong")
	}
}
