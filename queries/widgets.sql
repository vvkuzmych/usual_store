-- name: GetWidget :one
SELECT id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at
FROM widgets
WHERE id = $1;

-- name: ListWidgets :many
SELECT id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at
FROM widgets
ORDER BY name;

-- name: CreateWidget :one
INSERT INTO widgets (name, description, inventory_level, price, image, is_recurring, plan_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at;

-- name: UpdateWidget :one
UPDATE widgets
SET name = $2, description = $3, inventory_level = $4, price = $5, image = $6, is_recurring = $7, plan_id = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at;

-- name: DeleteWidget :exec
DELETE FROM widgets
WHERE id = $1;

-- name: GetWidgetsByRecurring :many
SELECT id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at
FROM widgets
WHERE is_recurring = $1
ORDER BY name;

-- name: UpdateWidgetInventory :one
UPDATE widgets
SET inventory_level = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at;

