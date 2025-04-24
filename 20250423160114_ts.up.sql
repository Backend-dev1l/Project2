-- Таблица пользователей 
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица балансов (создается только при первом зачислении)
CREATE TABLE IF NOT EXISTS balances (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    available DECIMAL(15, 2) NOT NULL DEFAULT 0.00 CHECK (available >= 0),
    reserved DECIMAL(15, 2) NOT NULL DEFAULT 0.00 CHECK (reserved >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица истории операций (для аудита)
CREATE TABLE IF NOT EXISTS operations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'reserve', 'revenue')),
    amount DECIMAL(15, 2) NOT NULL,
    order_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);