CREATE TABLE IF NOT EXISTS blocks (
  id SERIAL PRIMARY KEY,
  height INT NOT NULL,
  hash TEXT NOT NULL UNIQUE,
  prev_hash TEXT,
  timestamp BIGINT NOT NULL,
  validator TEXT,
  created_at TIMESTAMP DEFAULT now()
);
