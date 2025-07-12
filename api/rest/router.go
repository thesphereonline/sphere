package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/thesphereonline/sphere/db"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/blocks/", handleBlocks)
	mux.HandleFunc("/blocks", handleListBlocks)
	mux.HandleFunc("/tx/", handleTransaction)
	mux.HandleFunc("/portfolio/", handlePortfolio)

	return mux
}

func handleBlocks(w http.ResponseWriter, r *http.Request) {
	// URL: /blocks/{index}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid block path", http.StatusBadRequest)
		return
	}

	indexStr := parts[2]
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid block index", http.StatusBadRequest)
		return
	}

	block, err := db.GetBlockByIndex(index)
	if err != nil {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

func handleListBlocks(w http.ResponseWriter, r *http.Request) {
	blocks, err := db.ListBlocks()
	if err != nil {
		http.Error(w, "Failed to list blocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blocks)
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	// URL: /tx/{hash}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid transaction path", http.StatusBadRequest)
		return
	}

	hash := parts[2]
	tx, err := db.GetTransactionByHash(hash)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func handlePortfolio(w http.ResponseWriter, r *http.Request) {
	// URL: /portfolio/{address}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid portfolio path", http.StatusBadRequest)
		return
	}

	address := parts[2]
	portfolio, err := db.GetUserPortfolio(address)
	if err != nil {
		http.Error(w, "Failed to get portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portfolio)
}
