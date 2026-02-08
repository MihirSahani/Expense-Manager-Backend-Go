CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('cash', 'bank', 'credit_card', 'digital_wallet', 'investment')),
    currency VARCHAR(3) DEFAULT 'INR', -- ISO 4217 currency code
    current_balance DECIMAL(15, 2) DEFAULT 0.00,
    bank_name VARCHAR(100),
    account_number VARCHAR(50), -- last 4 digits only for security

    is_included_in_total BOOLEAN DEFAULT TRUE, -- exclude credit card debt from net worth
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);