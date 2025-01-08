package api

import (
	"encoding/json"
	"net/http"

	"sphere/core"
)

var blockchain = core.NewBlockchain()

func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain)
}

func StartServer() {
	http.HandleFunc("/blockchain", GetBlockchain)
	http.ListenAndServe(":8080", nil)
}
