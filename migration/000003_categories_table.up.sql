CREATE TABLE IF NOT EXISTS categories(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(255) NOT NULL CHECK (type IN ('income', 'expense')),
    color VARCHAR(7) NOT NULL DEFAULT '#000000',
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (name, user_id, type)
);

CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);