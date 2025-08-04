package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thesphereonline/sphere/node"
)

func blockRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getBlocksHandler)
	return r
}

func getBlocksHandler(w http.ResponseWriter, r *http.Request) {
	blocks := node.GlobalNode.Blockchain.Chain()
	json.NewEncoder(w).Encode(blocks)
}
