-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_banned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_is_banned ON users(is_banned);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Plans Table
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    price_toman BIGINT NOT NULL,
    traffic_bytes BIGINT NOT NULL,
    duration_days INTEGER NOT NULL
);

-- Extra Traffic Packages Table
CREATE TABLE IF NOT EXISTS extra_traffic_packages (
    id SERIAL PRIMARY KEY,
    price_toman BIGINT NOT NULL,
    traffic_bytes BIGINT NOT NULL
);

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('plan_purchase', 'extra_traffic_purchase')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    plan_id INTEGER REFERENCES plans(id),
    extra_traffic_id INTEGER REFERENCES extra_traffic_packages(id),
    price_toman INTEGER NOT NULL,
    authority VARCHAR(255) UNIQUE NOT NULL,
    ref_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_transactions_telegram_id ON transactions(telegram_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- Subscriptions Table (Unified Free and Paid)
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    plan_type VARCHAR(20) NOT NULL CHECK (plan_type IN ('free', 'paid')),
    traffic_limit_bytes BIGINT NOT NULL,
    traffic_used_bytes BIGINT DEFAULT 0,
    extra_traffic_limit_bytes BIGINT DEFAULT 0,
    transaction_id BIGINT REFERENCES transactions(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ DEFAULT (now() + interval '60 days'),
    is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_subscriptions_telegram_id ON subscriptions(telegram_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_active ON subscriptions(is_active);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_type ON subscriptions(plan_type);

-- File Logs Table
CREATE TABLE IF NOT EXISTS file_logs (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    channel_id BIGINT NOT NULL,
    file_id TEXT NOT NULL,
    file_name TEXT,
    file_type TEXT,
    file_size_byte BIGINT NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_file_logs_telegram_id ON file_logs(telegram_id);
CREATE INDEX IF NOT EXISTS idx_file_logs_created_at ON file_logs(created_at);

-- Admins Table
CREATE TABLE IF NOT EXISTS admins (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    telegram_id BIGINT REFERENCES users(telegram_id),
    password VARCHAR(255) NOT NULL
);
