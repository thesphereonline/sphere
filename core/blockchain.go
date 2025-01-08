package core

type Blockchain struct {
	Chain []Block
}

func NewBlockchain() *Blockchain {
	genesisBlock := Block{
		Index:     0,
		Timestamp: 0,
		PrevHash:  "",
		Hash:      "genesis",
		Validator: "Genesis Node",
	}
	return &Blockchain{Chain: []Block{genesisBlock}}
}

func (bc *Blockchain) AddBlock(newBlock Block) {
	if newBlock.PrevHash == bc.Chain[len(bc.Chain)-1].Hash {
		bc.Chain = append(bc.Chain, newBlock)
	}
}
