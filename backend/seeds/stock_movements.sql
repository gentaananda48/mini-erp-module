
INSERT INTO stock_movements (product_id, type, qty, note, created_at)
VALUES
    ((SELECT id FROM products WHERE sku = 'SKU-001'), 'in', 10, 'Initial stock', NOW() - INTERVAL '7 days'),
    ((SELECT id FROM products WHERE sku = 'SKU-002'), 'in', 25, 'Initial stock', NOW() - INTERVAL '7 days'),
    ((SELECT id FROM products WHERE sku = 'SKU-003'), 'in', 15, 'Initial stock', NOW() - INTERVAL '7 days'),
    ((SELECT id FROM products WHERE sku = 'SKU-001'), 'out', 2, 'Sample stock out', NOW() - INTERVAL '2 days'),
    ((SELECT id FROM products WHERE sku = 'SKU-002'), 'out', 1, 'Sample stock out', NOW() - INTERVAL '1 day');
