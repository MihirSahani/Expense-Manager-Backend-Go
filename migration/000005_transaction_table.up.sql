CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE SET NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE SET NULL,
    
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount >= 0),
    payee VARCHAR(255) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    transaction_date DATE DEFAULT NOW(),
    description TEXT,
    receipt_url VARCHAR(500) DEFAULT NULL,
    location VARCHAR(255) DEFAULT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date DESC);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date DESC);