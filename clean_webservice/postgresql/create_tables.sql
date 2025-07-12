-- Основная таблица заказов
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

-- Таблица доставки
CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL
);

-- Таблица платежей
CREATE TABLE IF NOT EXISTS payment (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

-- Таблица товаров (items)
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(50) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

-- Уникальный индекс для предотвращения дубликатов товаров в одном заказе
CREATE UNIQUE INDEX IF NOT EXISTS idx_items_order_chrt ON items(order_uid, chrt_id);
