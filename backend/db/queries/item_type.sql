-- name: GetAllItemTypes :many
SELECT * FROM item_types;

-- name: ListItemTypeOptions :many
SELECT id, name FROM item_types
ORDER BY name;

-- name: ListItemTypes :many
SELECT * FROM item_types
ORDER BY id
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountItemTypes :one
SELECT count(*) FROM item_types;

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
WHERE type_id = ANY(sqlc.arg('type_ids')::bigint[]);

-- name: AddItemTypeProperty :one
INSERT INTO item_type_properties (type_id, property_id, default_value, visibility) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateItemTypeProperty_DefaultValue :one
UPDATE item_type_properties SET default_value = $3 WHERE type_id = $1 AND property_id = $2 RETURNING *;

-- name: UpdateItemTypeProperty_Visibility :one
UPDATE item_type_properties SET visibility = $3 WHERE type_id = $1 AND property_id = $2 RETURNING *;

-- name: RemoveItemTypeProperty :execrows
DELETE FROM item_type_properties WHERE type_id = $1 AND property_id = $2;
