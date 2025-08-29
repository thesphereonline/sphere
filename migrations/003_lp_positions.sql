CREATE TABLE IF NOT EXISTS lp_positions (
  id SERIAL PRIMARY KEY,
  pool_id INT REFERENCES pools(id) ON DELETE CASCADE,
  owner TEXT NOT NULL,
  lp_amount TEXT NOT NULL DEFAULT '0',
  created_at TIMESTAMP DEFAULT now(),
  UNIQUE (pool_id, owner)
);
