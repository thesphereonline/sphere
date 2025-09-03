CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  block_id INT REFERENCES blocks(id) ON DELETE CASCADE,
  "from" TEXT,
  "to" TEXT,
  amount TEXT,
  fee TEXT,
  data TEXT,
  sig TEXT,
  created_at TIMESTAMP DEFAULT now()
);

-- index for quick lookups
CREATE INDEX IF NOT EXISTS idx_transactions_block_id ON transactions(block_id);
