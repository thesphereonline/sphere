package vm

const (
	OP_NOP        = 0x00
	OP_MINT       = 0x01 // Mint a token to an address
	OP_TRANSFER   = 0x02 // Transfer token from sender to another
	OP_BALANCEOF  = 0x03 // Push balance to stack
	OP_SWAP       = 0x04 // Execute AMM swap
	OP_MINT_NFT   = 0x05 // Mint NFT with metadata
	OP_STORE_META = 0x06 // Store NFT metadata (tokenId → URI)
)
