-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--bun:split

-- Create players table
CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

--bun:split

-- Create player sessions table
CREATE TABLE player_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id BIGINT REFERENCES players(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

--bun:split

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id BIGINT REFERENCES players(id) ON DELETE CASCADE,
    provider_id BIGINT UNIQUE,
    withdraw_provider_id BIGINT,
    amount VARCHAR(100) NOT NULL,
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('USD', 'EUR', 'KES')),
    status VARCHAR(12) NOT NULL CHECK (status IN ('PENDING', 'CONFIRMED', 'FAILED', 'FINAL', 'PROCESSING')),
    type VARCHAR(8) NOT NULL CHECK (type IN ('WITHDRAW', 'DEPOSIT', 'CANCEL')),
	attempts integer NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

--bun:split

CREATE INDEX idx_transactions_player_id ON transactions(player_id, status, created_at);

