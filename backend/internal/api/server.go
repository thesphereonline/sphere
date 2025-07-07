package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/thesphereonline/sphere/internal/core/amm"
	"github.com/thesphereonline/sphere/internal/core/mempool"
	"github.com/thesphereonline/sphere/internal/core/state"
	"github.com/thesphereonline/sphere/internal/core/transaction"
	"github.com/thesphereonline/sphere/internal/crypto"
)

type Server struct {
	Mempool *mempool.Mempool
	State   *state.State
}

func NewServer(m *mempool.Mempool, s *state.State) *Server {
	return &Server{Mempool: m, State: s}
}

func (s *Server) Start(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/submit-tx", s.submitTxHandler)
	mux.HandleFunc("/mint-nft", s.mintNFTHandler)
	mux.HandleFunc("/nfts", s.getNFTsHandler)
	mux.HandleFunc("/swap", s.swapHandler)
	mux.HandleFunc("/quote", s.getQuoteHandler)

	corsWrapped := WithCORS(mux)

	fmt.Println("🌐 REST API running on port", port)
	err := http.ListenAndServe(":"+port, corsWrapped)
	if err != nil {
		fmt.Printf("❌ REST API failed to start: %v\n", err)
	}
}

// --- Handlers ---
func (s *Server) submitTxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var tx transaction.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "invalid tx", 400)
		return
	}

	if !crypto.VerifySignature(tx.From, tx.Hash(), tx.Signature) {
		http.Error(w, "invalid signature", 400)
		return
	}

	fromAcc, err := s.State.GetAccount(tx.From)
	if err != nil {
		http.Error(w, "failed to load account", 500)
		return
	}

	if fromAcc.Nonce != tx.Nonce || fromAcc.Balance < tx.Amount {
		http.Error(w, "invalid nonce or insufficient balance", 400)
		return
	}

	if ok := s.Mempool.Add(tx); !ok {
		http.Error(w, "duplicate tx", 400)
		return
	}

	s.State.ValidateAndApplyTx(tx.From, tx.To, tx.Amount, tx.Nonce)
	w.WriteHeader(200)
	w.Write([]byte("✅ tx accepted"))
}

func (s *Server) mintNFTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var payload struct {
		Owner    string                 `json:"owner"`
		Name     string                 `json:"name"`
		ImageURL string                 `json:"image_url"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	err := s.State.DB.MintNFT(payload.Owner, payload.Name, payload.ImageURL, payload.Metadata)
	if err != nil {
		http.Error(w, "Mint failed: "+err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("✅ NFT minted"))
}

func (s *Server) getNFTsHandler(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	if owner == "" {
		http.Error(w, "Missing owner param", http.StatusBadRequest)
		return
	}

	nfts, err := s.State.DB.GetNFTsByOwner(owner)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Always return JSON array
	w.Header().Set("Content-Type", "application/json")
	if nfts == nil {
		nfts = []map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(nfts)
}

func (s *Server) swapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		From     string `json:"from"`
		TokenIn  string `json:"token_in"`
		TokenOut string `json:"token_out"`
		AmountIn uint64 `json:"amount_in"`
		Nonce    uint64 `json:"nonce"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	pool, err := s.State.DB.GetPool(req.TokenIn, req.TokenOut)
	if err != nil {
		http.Error(w, "Pool not found", 404)
		return
	}

	amountOut, err := amm.GetAmountOut(req.AmountIn, pool.ReserveA, pool.ReserveB)
	if err != nil {
		http.Error(w, "AMM error", 400)
		return
	}

	ok := s.State.ValidateAndApplyTx(req.From, "liquidity_pool", req.AmountIn, req.Nonce)
	if !ok {
		http.Error(w, "Transfer failed", 400)
		return
	}

	err = s.State.DB.UpdatePoolReserves(pool.TokenA, pool.TokenB, pool.ReserveA+req.AmountIn, pool.ReserveB-amountOut)
	if err != nil {
		http.Error(w, "Failed to update pool", 500)
		return
	}

	s.State.ValidateAndApplyTx("liquidity_pool", req.From, amountOut, 0)
	w.Write([]byte(fmt.Sprintf("✅ Swapped %d %s → %d %s", req.AmountIn, req.TokenIn, amountOut, req.TokenOut)))
}

func (s *Server) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	tokenIn := r.URL.Query().Get("token_in")
	tokenOut := r.URL.Query().Get("token_out")
	amountInStr := r.URL.Query().Get("amount_in")

	amountIn, err := strconv.ParseUint(amountInStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount_in", 400)
		return
	}

	pool, err := s.State.DB.GetPool(tokenIn, tokenOut)
	if err != nil {
		http.Error(w, "Pool not found", 404)
		return
	}

	amountOut, err := amm.GetAmountOut(amountIn, pool.ReserveA, pool.ReserveB)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]uint64{"amount_out": amountOut})
}
