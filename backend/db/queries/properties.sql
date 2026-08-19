-- name: GetAllProperties :many
SELECT * FROM properties;

-- name: ListPropertyOptions :many
SELECT id, name, value_type, default_value FROM properties
ORDER BY name;

-- name: ListProperties :many
SELECT * FROM properties
WHERE name ILIKE '%' || sqlc.arg('search') || '%'
ORDER BY id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountProperties :one
SELECT count(*) FROM properties
WHERE name ILIKE '%' || sqlc.arg('search') || '%';

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
