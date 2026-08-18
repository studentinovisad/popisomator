-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

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

-- name: CreateItems :many
INSERT INTO items (type_id) 
SELECT ($1) 
FROM generate_series(1, sqlc.arg(amount)::integer)
RETURNING *;

-- name: UpdateItem_Type :one
UPDATE items SET type_id = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItem_Consumption :one
UPDATE items SET consumption = $2 WHERE id = $1 RETURNING *;

-- name: DeleteItem :execrows
DELETE FROM items WHERE id = $1;

-- name: GetItemProperties :many
SELECT sqlc.embed(ip), itp.visibility, p.name AS property_name, it.derived_name_format FROM item_properties ip
JOIN items i ON ip.item_id = i.id
JOIN item_type_properties itp ON i.type_id = itp.type_id AND ip.property_id = itp.property_id 
JOIN item_types it ON i.type_id = it.id
JOIN properties p ON ip.property_id = p.id
WHERE ip.item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: AddItemProperty :one
INSERT INTO item_properties (item_id, property_id, property_value) VALUES ($1, $2, $3) RETURNING *;

-- name: AddItemPropertyBulk :many
INSERT INTO item_properties (item_id, property_id, property_value) 
SELECT 
    i.item_id, 
    sqlc.arg(property_id), 
    sqlc.arg(property_value)
FROM unnest(sqlc.arg(item_ids)::bigint[]) AS i(item_id)
RETURNING *;

-- name: UpdateItemProperty :one
UPDATE item_properties SET property_value = $3 WHERE item_id = $1 AND property_id = $2 RETURNING *;

-- name: RemoveItemProperty :execrows
DELETE FROM item_properties WHERE item_id = $1 AND property_id = $2;
