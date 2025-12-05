-- +goose Up
INSERT INTO orders (
    order_uuid,
    user_uuid,
    part_uuids,
    total_price,
    transaction_uuid,
    payment_method,
    status,
    created_at,
    updated_at
)
VALUES
-- 1. Новый заказ, ожидает оплаты
(
    gen_random_uuid(),
    gen_random_uuid(),
    ARRAY[gen_random_uuid(), gen_random_uuid()]::UUID[],
    1500.50,
    NULL,
    'CARD',
    'PENDING_PAYMENT',
    NOW(),
    NULL
),

-- 2. Оплаченный заказ
(
    gen_random_uuid(),
    gen_random_uuid(),
    ARRAY[gen_random_uuid()]::UUID[],
    299.99,
    gen_random_uuid(),
    'SBP',
    'PAID',
    NOW(),
    NOW()
),

-- 3. Отменённый заказ
(
    gen_random_uuid(),
    gen_random_uuid(),
    ARRAY[gen_random_uuid(), gen_random_uuid(), gen_random_uuid()]::UUID[],
    999.00,
    NULL,
    'INVESTOR_MONEY',
    'CANCELLED',
    NOW(),
    NOW()
);

-- +goose Down
DELETE FROM orders;