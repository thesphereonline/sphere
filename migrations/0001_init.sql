CREATE TABLE IF NOT EXISTS blocks (
    index BIGINT PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    prev_hash TEXT NOT NULL,
    hash TEXT NOT NULL,
    validator TEXT NOT NULL,
    state_root TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    hash TEXT PRIMARY KEY,
    block_index BIGINT REFERENCES blocks(index),
    from_addr TEXT NOT NULL,
    to_addr TEXT,
    nonce BIGINT NOT NULL,
    gas_limit BIGINT NOT NULL,
    data BYTEA NOT NULL,
    signature TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS balances (
    address TEXT NOT NULL,
    token TEXT NOT NULL DEFAULT 'SPHERE',
    balance BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (address, token)
);

CREATE TABLE IF NOT EXISTS nfts (
    token_id TEXT PRIMARY KEY,
    owner TEXT NOT NULL,
    metadata_uri TEXT
);

CREATE TABLE IF NOT EXISTS liquidity_pools (
    pool_id SERIAL PRIMARY KEY,
    token_a TEXT NOT NULL,
    token_b TEXT NOT NULL,
    reserve_a BIGINT NOT NULL,
    reserve_b BIGINT NOT NULL,
    total_shares BIGINT NOT NULL DEFAULT 0
);
