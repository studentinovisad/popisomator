-------- ITEMS

-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

-- name: GetAllItems :many
SELECT * FROM items;

-- name: ListItems :many
SELECT * FROM items
WHERE (sqlc.narg('type_id')::bigint IS NULL OR type_id = sqlc.narg('type_id'))
  AND (sqlc.narg('consumption')::consumption_status[] IS NULL OR consumption = ANY(sqlc.narg('consumption')::consumption_status[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to'))
ORDER BY
  CASE WHEN sqlc.arg('order_asc')::bool THEN created_at END ASC,
  CASE WHEN sqlc.arg('order_asc')::bool THEN id END ASC,
  CASE WHEN NOT sqlc.arg('order_asc')::bool THEN created_at END DESC,
  CASE WHEN NOT sqlc.arg('order_asc')::bool THEN id END DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountItems :one
SELECT count(*) FROM items
WHERE (sqlc.narg('type_id')::bigint IS NULL OR type_id = sqlc.narg('type_id'))
  AND (sqlc.narg('consumption')::consumption_status[] IS NULL OR consumption = ANY(sqlc.narg('consumption')::consumption_status[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to'));

-- name: CreateItem :one
INSERT INTO items (type_id) VALUES ($1) RETURNING *;

-- name: UpdateItemType :one
UPDATE items SET type_id = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItemConsumption :execrows
UPDATE items SET consumption = $2 WHERE id = $1;

-- name: DeleteItem :execrows
DELETE FROM items WHERE id = $1;

-------- PROPERTIES

-- name: GetAllProperties :many
SELECT * FROM properties;

-- name: GetPropertyByID :one
SELECT * FROM properties
WHERE id = $1 LIMIT 1;

-- name: CreateProperty :one
INSERT INTO properties (name, description, value_type, default_value) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePropertyName :exec
UPDATE properties SET name = $2 WHERE id = $1;

-- name: UpdatePropertyDescription :exec
UPDATE properties SET description = $2 WHERE id = $1;

-- name: UpdatePropertyDefaultValue :exec
UPDATE properties SET default_value = $2 WHERE id = $1;

-- name: DeleteProperty :execrows
DELETE FROM properties WHERE id = $1;


-------- ITEM PROPERTIES (PROPERTIES OF ITEMS)

-- name: GetItemProperties :many
SELECT * FROM item_properties
WHERE item_id = $1;

-- name: GetItemPropertiesForItems :many
SELECT * FROM item_properties
WHERE item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: AddItemProperty :one
INSERT INTO item_properties (item_id, property_id, property_value) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateItemProperty :one
UPDATE item_properties SET property_value = $3 WHERE item_id = $1 AND property_id = $2 RETURNING *;

-- name: RemoveItemProperty :execrows
DELETE FROM item_properties WHERE item_id = $1 AND property_id = $2;

-------- ITEM TYPES

-- name: GetAllItemTypes :many
SELECT * FROM item_types;

-- name: GetAllItemTypesWithProperties :many
SELECT
    sqlc.embed(t),
    tp.property_id,
    tp.default_value
FROM item_types t
LEFT JOIN item_type_properties tp ON tp.type_id = t.id
ORDER BY t.id;

-- name: GetItemTypeByID :one
SELECT * FROM item_types
WHERE id = $1 LIMIT 1;

-- name: CreateItemType :one
INSERT INTO item_types (name, description) VALUES ($1, $2) RETURNING *;

-- name: UpdateItemTypeName :one
UPDATE item_types SET name = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItemTypeDescription :one
UPDATE item_types SET description = $2 WHERE id = $1 RETURNING *;

-- name: DeleteItemType :execrows
DELETE FROM item_types WHERE id = $1;

-------- ITEM TYPE PROPERTIES (PROPERTIES TIED TO ITEM TYPES)

-- name: GetItemTypeProperties :many
SELECT * FROM item_type_properties
WHERE type_id = $1;

-- name: AddItemTypeProperty :one
INSERT INTO item_type_properties (type_id, property_id, default_value) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateItemTypeProperty :one
UPDATE item_type_properties SET default_value = $3 WHERE type_id = $1 AND property_id = $2 RETURNING *;

-- name: RemoveItemTypeProperty :execrows
DELETE FROM item_type_properties WHERE type_id = $1 AND property_id = $2;