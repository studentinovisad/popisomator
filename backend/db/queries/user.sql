-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE full_name ILIKE '%' || sqlc.arg(search)::text || '%'
  AND role = COALESCE(NULLIF(sqlc.arg(role_filter)::text, '')::user_role, role)
ORDER BY id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountUsers :one
SELECT count(*) FROM users
WHERE full_name ILIKE '%' || sqlc.arg(search)::text || '%'
  AND role = COALESCE(NULLIF(sqlc.arg(role_filter)::text, '')::user_role, role);

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, role) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateRole :one
UPDATE users SET role = $2 WHERE id = $1 RETURNING *;
