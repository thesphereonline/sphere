CREATE TABLE IF NOT EXISTS accounts (
    address TEXT PRIMARY KEY,
    balance BIGINT NOT NULL,
    nonce BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    hash TEXT PRIMARY KEY,
    sender TEXT NOT NULL,
    recipient TEXT NOT NULL,
    amount BIGINT NOT NULL,
    nonce BIGINT NOT NULL,
    signature TEXT,
    timestamp BIGINT
);

CREATE TABLE IF NOT EXISTS blocks (
    hash TEXT PRIMARY KEY,
    index BIGINT NOT NULL,
    timestamp BIGINT,
    prev_hash TEXT,
    validator TEXT,
    signature TEXT
);

CREATE TABLE IF NOT EXISTS block_transactions (
    block_hash TEXT REFERENCES blocks(hash),
    tx_hash TEXT REFERENCES transactions(hash),
    PRIMARY KEY (block_hash, tx_hash)
);

CREATE TABLE IF NOT EXISTS nfts (
  id SERIAL PRIMARY KEY,
  owner TEXT NOT NULL,
  name TEXT NOT NULL,
  image_url TEXT NOT NULL,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS liquidity_pools (
  id SERIAL PRIMARY KEY,
  token_a TEXT NOT NULL,
  token_b TEXT NOT NULL,
  reserve_a BIGINT DEFAULT 0,
  reserve_b BIGINT DEFAULT 0,
  lp_token TEXT NOT NULL,
  UNIQUE(token_a, token_b)
);

CREATE TABLE IF NOT EXISTS lp_balances (
  id SERIAL PRIMARY KEY,
  owner TEXT NOT NULL,
  lp_token TEXT NOT NULL,
  balance BIGINT DEFAULT 0,
  UNIQUE(owner, lp_token)
);

