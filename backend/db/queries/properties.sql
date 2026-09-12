-- Every property, or just the ones named. The filter is optional so one query serves both the full
-- catalogue and the audit path, which resolves a handful of names at once for an entry spanning
-- several properties - a reorder, or the initial list of a new item type.
-- name: GetProperties :many
SELECT * FROM properties
WHERE sqlc.narg('property_ids')::bigint[] IS NULL
   OR id = ANY(sqlc.narg('property_ids')::bigint[]);

-- name: ListPropertyOptions :many
SELECT id, name, value_type, default_value FROM properties
ORDER BY name;

-- name: ListProperties :many
SELECT * FROM properties
WHERE name ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%'
ORDER BY id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountProperties :one
SELECT count(*) FROM properties
WHERE name ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%';

-- name: GetPropertyByID :one
SELECT * FROM properties
WHERE id = $1 LIMIT 1;

-- name: CreateProperty :one
INSERT INTO properties (name, description, value_type, default_value) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProperty_Name :exec
UPDATE properties SET name = $2 WHERE id = $1;

-- name: UpdateProperty_Description :exec
UPDATE properties SET description = $2 WHERE id = $1;

-- name: UpdateProperty_DefaultValue :exec
UPDATE properties SET default_value = $2 WHERE id = $1;

-- name: DeleteProperty :execrows
DELETE FROM properties WHERE id = $1;
