INSERT INTO sales_orders (customer, status, total, created_at)
VALUES
    ('PT Demo Customer 1', 'CREATED', 37500000, NOW() - INTERVAL '2 days'),
    ('PT Demo Customer 2', 'CANCELLED', 3000000, NOW() - INTERVAL '1 day');

INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
VALUES
    (
        (SELECT id FROM sales_orders WHERE customer = 'PT Demo Customer 1'),
        'SKU-001',
        2,
        18000000
    ),
    (
        (SELECT id FROM sales_orders WHERE customer = 'PT Demo Customer 1'),
        'SKU-002',
        1,
        1500000
    );

INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
VALUES
    (
        (SELECT id FROM sales_orders WHERE customer = 'PT Demo Customer 2'),
        'SKU-002',
        2,
        1500000
    );
