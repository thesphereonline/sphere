package api

import (
	"encoding/json"
	"net/http"

	"github.com/thesphereonline/sphere/entities"

	"github.com/go-chi/chi/v5"
)

type WalletResponse struct {
	Address string `json:"address"`
}

func walletRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", createWalletHandler)
	return r
}

func createWalletHandler(w http.ResponseWriter, r *http.Request) {
	wallet, err := entities.NewWallet()
	if err != nil {
		http.Error(w, "failed to create wallet", http.StatusInternalServerError)
		return
	}

	resp := WalletResponse{
		Address: wallet.Address,
	}

	// TODO: Store wallet private keys securely or return to client
	// For now, return just the address (frontend manages keys)
	json.NewEncoder(w).Encode(resp)
}
