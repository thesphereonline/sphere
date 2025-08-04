package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Mount("/wallets", walletRouter())
		r.Mount("/transactions", transactionRouter())
		r.Mount("/blocks", blockRouter())
	})

	return r
}
