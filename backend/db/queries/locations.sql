-- name: CreateLocation :one
INSERT INTO locations (name, description, parent_id) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLocationByID :one
SELECT * FROM locations
WHERE id = $1 LIMIT 1;

-- name: ListLocationOptions :many
SELECT id, name, parent_id FROM locations
ORDER BY name;

-- name: CountLocations :one
SELECT count(*) FROM locations
WHERE name ILIKE '%' || escape_like_pattern(sqlc.arg('search')) || '%';

-- name: CheckLocationCycle :one
WITH RECURSIVE ancestors AS (
    SELECT id, parent_id FROM locations l WHERE l.id = cast(sqlc.arg(parent_id) as bigint)
    UNION ALL
    SELECT l.id, l.parent_id
    FROM locations l
    JOIN ancestors a ON l.id = a.parent_id
)
SELECT EXISTS (
  SELECT 1 FROM ancestors a WHERE a.id = cast(sqlc.arg(location_id) as bigint)
) AS would_create_cycle;

-- name: UpdateLocation_Name :exec
UPDATE locations SET name = $2 WHERE id = $1;

-- name: UpdateLocation_Description :exec
UPDATE locations SET description = $2 WHERE id = $1;

-- name: UpdateLocation_ParentID :exec
UPDATE locations SET parent_id = $2 WHERE id = $1;

-- name: DeleteLocation :execrows
DELETE FROM locations WHERE id = $1;