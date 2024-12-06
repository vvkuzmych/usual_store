-- Insert a record into the widgets table
INSERT INTO widgets (name, description, inventory_level, price, created_at, updated_at, image, is_recurring, plan_id)
VALUES ('Widget', 'A very nice widget.', 10, 1000, NOW(), NOW(), '/static/widget-1.png', false, '');

INSERT INTO widgets (name, description, inventory_level, price, created_at, updated_at, image, is_recurring, plan_id)
VALUES ('Golden Plan', 'Discount 30% for more than 3 subscriptions', 10, 3000, NOW(), NOW(), '', true, 'price_1PudhBRxsxaX9o1Hau9cfEqp');

