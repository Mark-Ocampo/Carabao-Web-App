package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mark-Ocampo/carabao-web-app/api/internal/account"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/auth"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/category"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/config"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/httpx"
	"github.com/Mark-Ocampo/carabao-web-app/api/internal/transaction"
)

func New(pool *pgxpool.Pool, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "db down"})
			return
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// auth
	authSvc := auth.NewService(auth.NewRepository(pool))
	authHandler := auth.NewHandler(authSvc, cfg.CookieSecure)
	r.Route("/auth", authHandler.Routes)

	// accounts
	accountHandler := account.NewHandler(account.NewService(account.NewRepository(pool)), authSvc)
	r.Route("/accounts", accountHandler.Routes)

	// categories
	categoryHandler := category.NewHandler(category.NewService(category.NewRepository(pool)), authSvc)
	r.Route("/categories", categoryHandler.Routes)

	// transactions
	transactionHandler := transaction.NewHandler(transaction.NewService(transaction.NewRepository(pool)), authSvc)
	r.Route("/transactions", transactionHandler.Routes)

	return r
}
