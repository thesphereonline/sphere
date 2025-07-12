package types

type State struct {
	Balances map[string]uint64            // address → balance
	Tokens   map[string]map[string]uint64 // tokenAddr → (address → balance)
	NFTs     map[string]string            // tokenId → owner
	Metadata map[string]string            // tokenId → URI
}
