-- name: GetAllItemTypes :many
SELECT * FROM item_types;

-- name: ListItemTypeOptions :many
SELECT id, name FROM item_types
ORDER BY name;

-- name: ListItemTypes :many
SELECT * FROM item_types
WHERE name ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%'
ORDER BY id
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountItemTypes :one
SELECT count(*) FROM item_types
WHERE name ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%';

-- name: GetItemTypeByID :one
SELECT * FROM item_types
WHERE id = $1 LIMIT 1;

-- name: CreateItemType :one
INSERT INTO item_types (name, description, derived_name_format) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateItemType_Name :one
UPDATE item_types SET name = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItemType_Description :one
UPDATE item_types SET description = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItemType_DerivedNameFormat :one
UPDATE item_types SET derived_name_format = $2 WHERE id = $1 RETURNING *;

-- name: DeleteItemType :execrows
DELETE FROM item_types WHERE id = $1;

-- name: GetItemTypeProperties :many
SELECT * FROM item_type_properties
WHERE type_id = ANY(sqlc.arg('type_ids')::bigint[])
ORDER BY type_id, position;

-- name: AddItemTypeProperty :one
WITH locked_type AS (
  SELECT id FROM item_types WHERE id = $1 FOR UPDATE
), next_position AS (
  SELECT COALESCE(max(position), -1) + 1 AS position
  FROM item_type_properties
  WHERE type_id = $1
)
INSERT INTO item_type_properties (type_id, property_id, default_value, visibility, position)
SELECT $1, $2, $3, $4, next_position.position
FROM locked_type CROSS JOIN next_position
RETURNING *;

-- name: UpdateItemTypeProperty_DefaultValue :one
UPDATE item_type_properties SET default_value = $3 WHERE type_id = $1 AND property_id = $2 RETURNING *;

-- name: UpdateItemTypeProperty_Visibility :one
UPDATE item_type_properties SET visibility = $3 WHERE type_id = $1 AND property_id = $2 RETURNING *;

-- name: RemoveItemTypeProperty :execrows
DELETE FROM item_type_properties WHERE type_id = $1 AND property_id = $2;
