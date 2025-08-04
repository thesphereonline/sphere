-- +migrate Up
SET search_path TO public;

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    public_key TEXT NOT NULL,
    address TEXT UNIQUE NOT NULL,
    nonce BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hash TEXT UNIQUE NOT NULL,
    height BIGINT NOT NULL,
    parent_hash TEXT,
    proposer TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hash TEXT UNIQUE NOT NULL,
    sender TEXT NOT NULL,
    recipient TEXT,
    amount BIGINT,
    gas_fee BIGINT DEFAULT 0,
    payload JSONB,
    block_hash TEXT,
    signature TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    decimals INT DEFAULT 18,
    total_supply BIGINT DEFAULT 0,
    owner TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE nfts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_id TEXT NOT NULL,
    owner TEXT NOT NULL,
    metadata_uri TEXT,
    minted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stakes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staker TEXT NOT NULL,
    validator TEXT NOT NULL,
    amount BIGINT NOT NULL,
    delegated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE governance_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    proposal_id TEXT NOT NULL,
    voter TEXT NOT NULL,
    vote_choice TEXT,
    weight BIGINT DEFAULT 1,
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
