-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

-- name: ListItems :many
SELECT items.* FROM items
JOIN item_types ON item_types.id = items.type_id
WHERE (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
  AND (sqlc.narg('consumption')::consumption_status[] IS NULL OR items.consumption = ANY(sqlc.narg('consumption')::consumption_status[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR items.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR items.created_at <= sqlc.narg('created_to'))
  AND (
    sqlc.arg('search')::text = ''
    OR item_types.derived_name_format ILIKE '%' || escape_like_pattern(sqlc.arg('search')::text) || '%'
    OR EXISTS (
      SELECT 1 FROM item_properties
      JOIN properties ON properties.id = item_properties.property_id
      WHERE item_properties.item_id = items.id
        AND item_types.derived_name_format LIKE '%{' || escape_like_pattern(properties.name) || '}%'
        AND format_property_value(item_properties.property_value, properties.value_type)
          ILIKE '%' || escape_like_pattern(sqlc.arg('search')::text) || '%'
    )
    OR render_item_derived_name(items.id, item_types.derived_name_format)
      ILIKE '%' || replace(escape_like_pattern(trim(sqlc.arg('search')::text)), ' ', '%') || '%'
  )
  AND NOT EXISTS (
    SELECT 1
    FROM ROWS FROM (
      unnest(sqlc.arg('property_ids')::bigint[]),
      unnest(sqlc.arg('property_values')::jsonb[])
    ) AS filters(property_id, property_value)
    WHERE NOT EXISTS (
      SELECT 1
      FROM item_properties
      WHERE item_properties.item_id = items.id
        AND item_properties.property_id = filters.property_id
        AND item_properties.property_value = filters.property_value
    )
  )
ORDER BY
  CASE WHEN sqlc.arg('order_asc')::bool THEN items.created_at END ASC,
  CASE WHEN sqlc.arg('order_asc')::bool THEN items.id END ASC,
  CASE WHEN NOT sqlc.arg('order_asc')::bool THEN items.created_at END DESC,
  CASE WHEN NOT sqlc.arg('order_asc')::bool THEN items.id END DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountItems :one
SELECT count(*) FROM items
JOIN item_types ON item_types.id = items.type_id
WHERE (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
  AND (sqlc.narg('consumption')::consumption_status[] IS NULL OR items.consumption = ANY(sqlc.narg('consumption')::consumption_status[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR items.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR items.created_at <= sqlc.narg('created_to'))
  AND (
    sqlc.arg('search')::text = ''
    OR item_types.derived_name_format ILIKE '%' || escape_like_pattern(sqlc.arg('search')::text) || '%'
    OR EXISTS (
      SELECT 1 FROM item_properties
      JOIN properties ON properties.id = item_properties.property_id
      WHERE item_properties.item_id = items.id
        AND item_types.derived_name_format LIKE '%{' || escape_like_pattern(properties.name) || '}%'
        AND format_property_value(item_properties.property_value, properties.value_type)
          ILIKE '%' || escape_like_pattern(sqlc.arg('search')::text) || '%'
    )
    OR render_item_derived_name(items.id, item_types.derived_name_format)
      ILIKE '%' || replace(escape_like_pattern(trim(sqlc.arg('search')::text)), ' ', '%') || '%'
  )
  AND NOT EXISTS (
    SELECT 1
    FROM ROWS FROM (
      unnest(sqlc.arg('property_ids')::bigint[]),
      unnest(sqlc.arg('property_values')::jsonb[])
    ) AS filters(property_id, property_value)
    WHERE NOT EXISTS (
      SELECT 1
      FROM item_properties
      WHERE item_properties.item_id = items.id
        AND item_properties.property_id = filters.property_id
        AND item_properties.property_value = filters.property_value
    )
  );

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
SELECT sqlc.embed(ip), itp.visibility, p.name AS property_name, p.value_type AS property_type FROM item_properties ip
JOIN items i ON ip.item_id = i.id
JOIN item_type_properties itp ON i.type_id = itp.type_id AND ip.property_id = itp.property_id
JOIN item_types it ON i.type_id = it.id
JOIN properties p ON ip.property_id = p.id
WHERE ip.item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: ListItemTypeFilterableProperties :many
SELECT
  item_type_properties.property_id,
  count(DISTINCT item_properties.property_value) AS value_count
FROM item_type_properties
LEFT JOIN items ON items.type_id = item_type_properties.type_id
LEFT JOIN item_properties
  ON item_properties.item_id = items.id
  AND item_properties.property_id = item_type_properties.property_id
WHERE item_type_properties.type_id = $1
GROUP BY item_type_properties.property_id
ORDER BY item_type_properties.property_id;

-- name: ListItemTypePropertyValues :many
SELECT DISTINCT item_properties.property_value AS value
FROM item_properties
JOIN items ON items.id = item_properties.item_id
JOIN item_type_properties
  ON item_type_properties.type_id = items.type_id
  AND item_type_properties.property_id = item_properties.property_id
JOIN properties ON properties.id = item_properties.property_id
WHERE items.type_id = sqlc.arg('type_id')
  AND item_properties.property_id = sqlc.arg('property_id')
  AND format_property_value(item_properties.property_value, properties.value_type)
    ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%'
ORDER BY value
LIMIT sqlc.arg('limit_val');

-- name: GetItemsDerivedNames :many
SELECT items.id AS item_id, render_item_derived_name(items.id, item_types.derived_name_format) AS derived_name
FROM items
JOIN item_types ON item_types.id = items.type_id
WHERE items.id = ANY(sqlc.arg('item_ids')::bigint[]);

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
