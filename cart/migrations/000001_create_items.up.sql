CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    sku BIGINT UNIQUE NOT NULL,
    count INT NOT NULL,
    CONSTRAINT unique_user_sku_pair UNIQUE (user_id, sku)
);
