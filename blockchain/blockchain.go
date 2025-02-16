package blockchain

import (
	"fmt"
	"your-username/blockchain/block"
	"your-username/blockchain/storage"
	"your-username/blockchain/transaction"
	"your-username/blockchain/utils"
)

type Blockchain struct {
	Blocks       []*block.Block
	PendingTxs   []*transaction.Transaction
	Storage      *storage.Storage
	Difficulty   int
	MiningReward float64
}

// CreateBlockchain initializes a new blockchain with genesis block
func CreateBlockchain(dbFile string) *Blockchain {
	storage := storage.NewStorage(dbFile)
	blocks, err := storage.LoadChain()
	if err != nil {
		panic(err)
	}

	if len(blocks) == 0 {
		genesis := block.CreateBlock([]*transaction.Transaction{}, "", 4)
		blocks = append(blocks, genesis)
	}

	return &Blockchain{
		Blocks:       blocks,
		Storage:      storage,
		Difficulty:   4,
		MiningReward: 100,
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []*transaction.Transaction) *block.Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := block.CreateBlock(transactions, prevBlock.Hash, bc.Difficulty)
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.Storage.SaveChain(bc.Blocks)
	return newBlock
}

func (bc *Blockchain) AddTransaction(tx *transaction.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	// Validate transaction
	if tx.Amount <= 0 {
		return fmt.Errorf("invalid transaction amount")
	}

	// Add to pending transactions
	bc.PendingTxs = append(bc.PendingTxs, tx)
	return nil
}

func (bc *Blockchain) GetBalance(address string) float64 {
	balance := 0.0

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.From == address {
				balance -= tx.Amount
			}
			if tx.To == address {
				balance += tx.Amount
			}
		}
	}

	return balance
}

func (bc *Blockchain) ValidateChain() error {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// Validate block hash
		if currentBlock.PrevBlockHash != previousBlock.Hash {
			return fmt.Errorf("invalid block chain: block %d has invalid previous hash", i)
		}

		// Validate block hash calculation
		if currentBlock.CalculateHash() != currentBlock.Hash {
			return fmt.Errorf("invalid block chain: block %d has invalid hash", i)
		}
	}
	return nil
}

func (bc *Blockchain) ValidateAndAddBlock(block *block.Block) error {
	// Validate block
	if block == nil {
		return utils.NewError("block cannot be nil", nil)
	}

	// Check if block already exists
	for _, b := range bc.Blocks {
		if b.Hash == block.Hash {
			return utils.NewError("block already exists", nil)
		}
	}

	// Validate previous hash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	if block.PrevBlockHash != lastBlock.Hash {
		return utils.NewError("invalid previous hash", nil)
	}

	// Validate block hash
	if block.CalculateHash() != block.Hash {
		return utils.NewError("invalid block hash", nil)
	}

	// Add block
	bc.Blocks = append(bc.Blocks, block)
	return nil
}
