package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"your-username/blockchain/block"
	"your-username/blockchain/blockchain"
	"your-username/blockchain/transaction"
	"your-username/blockchain/utils"
)

type Server struct {
	Blockchain    *blockchain.Blockchain
	Peers         map[string]bool
	MinerAddress  string
	mu            sync.RWMutex
	server        *http.Server
	isMining      bool
	miningStopped chan struct{}
}

func NewServer(bc *blockchain.Blockchain) *Server {
	return &Server{
		Blockchain:    bc,
		Peers:         make(map[string]bool),
		miningStopped: make(chan struct{}),
	}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Blockchain endpoints
	mux.HandleFunc("/blocks", s.handleBlocks)
	mux.HandleFunc("/mine", s.handleMine)
	mux.HandleFunc("/transaction", s.handleTransaction)
	mux.HandleFunc("/peers", s.handlePeers)
	mux.HandleFunc("/sync", s.handleSync)
	mux.HandleFunc("/status", s.handleStatus)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start mining loop
	go s.miningLoop()

	// Start peer synchronization
	go s.syncLoop()

	return s.server.ListenAndServe()
}

func (s *Server) Stop() {
	s.isMining = false
	<-s.miningStopped
	if s.server != nil {
		s.server.Close()
	}
}

func (s *Server) AddPeer(address string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if address == "" {
		return
	}

	if _, exists := s.Peers[address]; !exists {
		s.Peers[address] = true
		log.Printf("Added new peer: %s", address)
		go s.syncWithPeer(address)
	}
}

func (s *Server) miningLoop() {
	s.isMining = true
	for s.isMining {
		if len(s.Blockchain.PendingTxs) > 0 {
			newBlock := s.Blockchain.AddBlock(s.Blockchain.PendingTxs)
			s.Blockchain.PendingTxs = nil
			s.broadcastNewBlock(newBlock)
		}
		time.Sleep(10 * time.Second)
	}
	s.miningStopped <- struct{}{}
}

func (s *Server) syncLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.mu.RLock()
		peers := make([]string, 0, len(s.Peers))
		for peer := range s.Peers {
			peers = append(peers, peer)
		}
		s.mu.RUnlock()

		for _, peer := range peers {
			go s.syncWithPeer(peer)
		}
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := struct {
		Blocks       int    `json:"blocks"`
		Peers        int    `json:"peers"`
		PendingTxs   int    `json:"pendingTransactions"`
		MinerAddress string `json:"minerAddress"`
	}{
		Blocks:       len(s.Blockchain.Blocks),
		Peers:        len(s.Peers),
		PendingTxs:   len(s.Blockchain.PendingTxs),
		MinerAddress: s.MinerAddress,
	}

	json.NewEncoder(w).Encode(status)
}

func (s *Server) handleBlocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(s.Blockchain.Blocks)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newBlock := s.Blockchain.AddBlock(s.Blockchain.PendingTxs)
	s.Blockchain.PendingTxs = nil

	// Broadcast to peers
	s.broadcastNewBlock(newBlock)

	json.NewEncoder(w).Encode(newBlock)
}

func (s *Server) handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTx := transaction.NewTransaction(tx.From, tx.To, tx.Amount)
	if err := s.Blockchain.AddTransaction(newTx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) broadcastNewBlock(block *block.Block) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for peer := range s.Peers {
		// Implement peer notification logic
		go s.notifyPeer(peer, block)
	}
}

func (s *Server) syncWithPeer(peer string) {
	resp, err := http.Get(fmt.Sprintf("http://%s/blocks", peer))
	if err != nil {
		log.Printf("Failed to sync with peer %s: %v", peer, err)
		return
	}
	defer resp.Body.Close()

	var blocks []*block.Block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		log.Printf("Failed to decode blocks from peer %s: %v", peer, err)
		return
	}

	// Compare and update blockchain if necessary
	if len(blocks) > len(s.Blockchain.Blocks) {
		s.Blockchain.Blocks = blocks
	}
}

func (s *Server) notifyPeer(peer string, block *block.Block) error {
	blockData, err := json.Marshal(block)
	if err != nil {
		return utils.NewError("failed to marshal block", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/sync", peer), "application/json", bytes.NewBuffer(blockData))
	if err != nil {
		return utils.NewError("failed to notify peer", err)
	}
	defer resp.Body.Close()
	return nil
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var block block.Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate and add block
	if err := s.Blockchain.ValidateAndAddBlock(&block); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlePeers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.mu.RLock()
		peers := make([]string, 0, len(s.Peers))
		for peer := range s.Peers {
			peers = append(peers, peer)
		}
		s.mu.RUnlock()
		json.NewEncoder(w).Encode(peers)

	case "POST":
		var peer struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
			http.Error(w, "Invalid peer address", http.StatusBadRequest)
			return
		}

		s.AddPeer(peer.Address)
		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
