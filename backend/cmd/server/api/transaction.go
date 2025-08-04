package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thesphereonline/sphere/entities"
	"github.com/thesphereonline/sphere/node"
)

type CreateTransactionRequest struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
	// Signature can be added here or signed client-side
}

func transactionRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", createTransactionHandler)
	return r
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx := entities.Transaction{
		ID:        "", // generate UUID or hash
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}

	// TODO: Verify signature before accepting

	// Add transaction to mempool
	node.GlobalNode.TransactionPool.Add(tx)

	json.NewEncoder(w).Encode(tx)
}
