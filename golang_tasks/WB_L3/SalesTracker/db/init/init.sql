CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income','expense')),
    amount NUMERIC(14,2) NOT NULL CHECK (amount >= 0),
    category VARCHAR(100) DEFAULT 'uncategorized',
    note TEXT,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
