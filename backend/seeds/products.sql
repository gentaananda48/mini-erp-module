
INSERT INTO products (sku, name, price, stock, is_active)
VALUES
    ('SKU-001', 'MacBook Air M3', 18000000, 10, TRUE),
    ('SKU-002', 'Logitech MX Master 3S', 1500000, 25, TRUE),
    ('SKU-003', 'Keychron K2 Mechanical Keyboard', 1200000, 15, TRUE),
    ('SKU-004', 'Dell UltraSharp Monitor 27', 6500000, 8, TRUE),
    ('SKU-005', 'Samsung T7 SSD 1TB', 1800000, 20, TRUE),
    ('SKU-006', 'Logitech C920 Webcam', 950000, 12, TRUE),
    ('SKU-007', 'Anker USB-C Hub', 750000, 30, TRUE),
    ('SKU-008', 'Sony WH-1000XM5', 4500000, 5, TRUE),
    ('SKU-009', 'USB-C Cable 2M', 150000, 50, TRUE),
    ('SKU-010', 'Laptop Stand Aluminum', 350000, 0, FALSE),
    ('SKU-011', 'Figure Model Gundam Seed Freedom', 450000, 20, TRUE),
    ('SKU-012', 'Vacuum Insulate Bottle', 75000, 10, TRUE),
    ('SKU-013', 'Rubix 3x3', 65000, 19, TRUE),
    ('SKU-014', 'Snowman Drawing Pen', 5000, 20, TRUE)
ON CONFLICT (sku) DO NOTHING;
