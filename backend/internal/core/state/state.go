package state

import (
	"sync"

	"github.com/thesphereonline/sphere/internal/db"
)

type Account struct {
	Balance uint64
	Nonce   uint64
}

type State struct {
	DB   *db.Postgres
	lock sync.Mutex
}

// NewState creates a new State with a connected Postgres instance.
func NewState(pg *db.Postgres) *State {
	return &State{DB: pg}
}

// GetAccount retrieves account info from the database.
// If the account does not exist, creates with initial faucet balance and nonce=0.
func (s *State) GetAccount(addr string) (*Account, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	balance, nonce, err := s.DB.GetAccount(addr)
	if err != nil {
		// Account probably does not exist, create it with initial balance
		initialBalance := uint64(1000)
		initialNonce := uint64(0)
		err := s.DB.UpsertAccount(addr, initialBalance, initialNonce)
		if err != nil {
			return nil, err
		}
		return &Account{Balance: initialBalance, Nonce: initialNonce}, nil
	}
	return &Account{Balance: balance, Nonce: nonce}, nil
}

// ValidateAndApplyTx verifies transaction validity and updates balances atomically.
// Returns true if transaction applied successfully, false otherwise.
func (s *State) ValidateAndApplyTx(from, to string, amount, nonce uint64) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Fetch sender account
	fromBal, fromNonce, err := s.DB.GetAccount(from)
	if err != nil {
		return false // sender account does not exist
	}

	// Basic validation: sufficient balance and nonce must match
	if fromBal < amount || fromNonce != nonce {
		return false
	}

	// Fetch or create recipient account
	toBal, _, err := s.DB.GetAccount(to)
	if err != nil {
		// Recipient does not exist — create with 0 balance and nonce 0
		err = s.DB.UpsertAccount(to, 0, 0)
		if err != nil {
			return false
		}
		toBal = 0
	}

	// Perform atomic update
	err = s.DB.UpsertAccount(from, fromBal-amount, nonce+1)
	if err != nil {
		return false
	}
	err = s.DB.UpsertAccount(to, toBal+amount, 0)
	if err != nil {
		return false
	}

	return true
}
