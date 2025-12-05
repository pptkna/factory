-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL,
    total_price REAL NOT NUll,
    transaction_uuid UUID,
    payment_method VARCHAR(32),
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE orders;