CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) CHECK (type IN ('Stock', 'Gold', 'Bond', 'Savings', 'Crypto')) NOT NULL,
    ticker VARCHAR(20), -- Optional, only for stocks/crypto
    price NUMERIC(18,2), -- Can be NULL for non-priced assets
    amount NUMERIC(18,4) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(10) DEFAULT 'USD',
    interest_rate NUMERIC(5,2), -- Optional, only for Savings/Bonds
    compounding_frequency VARCHAR(20) CHECK (compounding_frequency IN ('daily', 'monthly', 'quarterly', 'annually')),
    interest_start DATE, -- Only for interest-based assets
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE asset_returns (
    id SERIAL PRIMARY KEY,
    asset_id INT REFERENCES assets(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    returns NUMERIC(18,4) NOT NULL
);
