package types

type Transaction struct {
	Hash       string `json:"hash"`
	From       string `json:"from"`
	To         string `json:"to"`
	Nonce      uint64 `json:"nonce"`
	GasLimit   uint64 `json:"gas_limit"`
	Data       []byte `json:"data"`
	Signature  string `json:"signature"`
	BlockIndex uint64 `json:"block_index"` // ADD THIS
}

type Block struct {
	Index     uint64
	Timestamp int64
	PrevHash  string
	Hash      string
	Validator string
	Txs       []Transaction
	StateRoot string
	Signature string
}

type Portfolio struct {
	Address     string                  `json:"address"`
	Balances    []TokenBalance          `json:"balances"`
	NFTs        []NFT                   `json:"nfts"`
	LPPositions []LiquidityPoolPosition `json:"lp_positions"`
}

type TokenBalance struct {
	Token   string `json:"token"`
	Balance uint64 `json:"balance"`
}

type NFT struct {
	TokenID     string `json:"token_id"`
	MetadataURI string `json:"metadata_uri"`
	Owner       string `json:"owner"`
}

type LiquidityPoolPosition struct {
	PoolID     uint64 `json:"pool_id"`
	TokenA     string `json:"token_a"`
	TokenB     string `json:"token_b"`
	Shares     uint64 `json:"shares"`
	ShareValue uint64 `json:"share_value"` // maybe USD or token units
}
