package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thesphereonline/sphere/internal/core/block"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(connString string) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	fmt.Println("📦 Connected to PostgreSQL")
	return &Postgres{Pool: pool}, nil
}

func (pg *Postgres) GetAccount(address string) (uint64, uint64, error) {
	var balance, nonce uint64
	err := pg.Pool.QueryRow(context.Background(),
		"SELECT balance, nonce FROM accounts WHERE address=$1", address).Scan(&balance, &nonce)
	if err != nil {
		return 0, 0, err
	}
	return balance, nonce, nil
}

func (pg *Postgres) UpsertAccount(address string, balance, nonce uint64) error {
	_, err := pg.Pool.Exec(context.Background(),
		`INSERT INTO accounts (address, balance, nonce)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (address)
		 DO UPDATE SET balance = $2, nonce = $3`,
		address, balance, nonce)
	return err
}

// SaveBlock persists a block and its transactions
func (pg *Postgres) SaveBlock(b *block.Block) error {
	tx, err := pg.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`INSERT INTO blocks (hash, index, timestamp, prev_hash, validator, signature)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (hash) DO NOTHING`,
		b.Hash, b.Index, b.Timestamp, b.PrevHash, b.Validator, b.Signature)
	if err != nil {
		return err
	}

	for _, t := range b.Txs {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO transactions (hash, sender, recipient, amount, nonce, signature, timestamp)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (hash) DO NOTHING`,
			t.Hash(), t.From, t.To, t.Amount, t.Nonce, t.Signature, t.Timestamp)
		if err != nil {
			return err
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO block_transactions (block_hash, tx_hash)
			 VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`,
			b.Hash, t.Hash())
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (pg *Postgres) LoadBlocks() ([]*block.Block, error) {
	rows, err := pg.Pool.Query(context.Background(), `SELECT hash, index, timestamp, prev_hash, validator, signature FROM blocks ORDER BY index ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*block.Block
	for rows.Next() {
		var b block.Block
		err := rows.Scan(&b.Hash, &b.Index, &b.Timestamp, &b.PrevHash, &b.Validator, &b.Signature)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &b)
	}
	return blocks, nil
}

func (p *Postgres) MintNFT(owner, name, imageURL string, metadata map[string]interface{}) error {
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	_, err = p.Pool.Exec(context.Background(), `
		INSERT INTO nfts (owner, name, image_url, metadata)
		VALUES ($1, $2, $3, $4)
	`, owner, name, imageURL, metaJSON)
	return err
}

type NFT struct {
	ID       int                    `json:"id"`
	Owner    string                 `json:"owner"`
	Name     string                 `json:"name"`
	ImageURL string                 `json:"image_url"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (p *Postgres) GetNFTsByOwner(owner string) ([]map[string]interface{}, error) {
	rows, err := p.Pool.Query(context.Background(),
		"SELECT name, image_url, metadata FROM nfts WHERE owner=$1", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var name, imageURL string
		var metadataJSON []byte

		err := rows.Scan(&name, &imageURL, &metadataJSON)
		if err != nil {
			return nil, err
		}

		var metadata map[string]interface{}
		json.Unmarshal(metadataJSON, &metadata)

		result = append(result, map[string]interface{}{
			"name":      name,
			"image_url": imageURL,
			"metadata":  metadata,
		})
	}
	return result, nil
}

type Pool struct {
	TokenA   string
	TokenB   string
	ReserveA uint64
	ReserveB uint64
	LPToken  string
}

func (p *Postgres) GetPool(tokenA, tokenB string) (*Pool, error) {
	row := p.Pool.QueryRow(context.Background(), `
		SELECT token_a, token_b, reserve_a, reserve_b, lp_token
		FROM liquidity_pools
		WHERE token_a=$1 AND token_b=$2
	`, tokenA, tokenB)

	var pool Pool
	err := row.Scan(&pool.TokenA, &pool.TokenB, &pool.ReserveA, &pool.ReserveB, &pool.LPToken)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

func (p *Postgres) UpdatePoolReserves(tokenA, tokenB string, newA, newB uint64) error {
	_, err := p.Pool.Exec(context.Background(), `
		UPDATE liquidity_pools
		SET reserve_a = $3, reserve_b = $4
		WHERE token_a=$1 AND token_b=$2
	`, tokenA, tokenB, newA, newB)
	return err
}
