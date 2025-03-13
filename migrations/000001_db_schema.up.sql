BEGIN;

CREATE TYPE transaction_type AS ENUM ('deposit', 'withdraw');

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(15) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(16) NOT NULL UNIQUE,
    phone VARCHAR(13) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE balances (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    balance DECIMAL(15, 2) NOT NULL CHECK >= 0 DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE balances
ADD CONSTRAINT check_balance CHECK (balance >= 0)

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    type transaction_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;