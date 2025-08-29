CREATE TABLE IF NOT EXISTS validators (
  address TEXT PRIMARY KEY,
  stake TEXT NOT NULL DEFAULT '0',
  commission_bps INT DEFAULT 1000, -- 10% default
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS delegations (
  id SERIAL PRIMARY KEY,
  delegator TEXT NOT NULL,
  validator TEXT NOT NULL REFERENCES validators(address),
  amount TEXT NOT NULL DEFAULT '0',
  created_at TIMESTAMP DEFAULT now()
);
