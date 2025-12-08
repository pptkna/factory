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
    '8f6bd941-9b0a-4fcd-9c43-0e8ece9a41e5'::UUID,  -- order UUID
    '4e12b9f1-ec1a-42c7-974b-61644d63f2e1'::UUID,  -- user UUID
    ARRAY[
        'c2a4e5cd-8aa1-4dd8-ab19-5cdf25af047d'::UUID,
        '42d9bd89-6023-4d95-af31-c1d86bb39b43'::UUID
    ]::UUID[],
    1500.50,
    NULL,
    'CARD',
    'PENDING_PAYMENT',
    NOW(),
    NULL
),

-- 2. Оплаченный заказ
(
    'd3e9e4a7-8f51-4fc7-af0f-92c6768a3dd8'::UUID,
    'b90c6d18-9746-4aaf-bb7b-77ee71e2c524'::UUID,
    ARRAY[
        '19b5cae8-54d3-4fc8-97fa-bab9d0718d80'::UUID
    ]::UUID[],
    299.99,
    'a2f7f078-e664-4bba-a740-2485c4cebd61'::UUID,
    'SBP',
    'PAID',
    NOW(),
    NOW()
),

-- 3. Отменённый заказ
(
    '1c78c089-44d6-4e2f-8dcb-6a37a4a4d4c5'::UUID,
    'eaa2195a-65af-4d2b-8a5d-8e50f7af5d93'::UUID,
    ARRAY[
        'c2a4e5cd-8aa1-4dd8-ab19-5cdf25af047d'::UUID,
        '42d9bd89-6023-4d95-af31-c1d86bb39b43'::UUID,
        '19b5cae8-54d3-4fc8-97fa-bab9d0718d80'::UUID
    ]::UUID[],
    999.00,
    NULL,
    'INVESTOR_MONEY',
    'CANCELLED',
    NOW(),
    NOW()
);

-- +goose Down
DELETE FROM orders;