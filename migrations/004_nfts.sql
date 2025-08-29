CREATE TABLE IF NOT EXISTS nfts (
  id SERIAL PRIMARY KEY,
  token_id TEXT UNIQUE NOT NULL,
  owner TEXT NOT NULL,
  metadata JSONB DEFAULT '{}',
  royalties_bps INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT now()
);
