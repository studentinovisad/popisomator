-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users u
WHERE u.full_name ILIKE '%' || escape_like_pattern(sqlc.arg(search)::text) || '%'
  AND (sqlc.narg('role')::user_role IS NULL OR u.role = sqlc.narg('role')::user_role)
  AND (sqlc.narg('status')::user_status IS NULL OR u.status = sqlc.narg('status')::user_status)
ORDER BY id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountUsers :one
SELECT count(*) FROM users u
WHERE u.full_name ILIKE '%' || escape_like_pattern(sqlc.arg(search)::text) || '%'
  AND (sqlc.narg('role')::user_role IS NULL OR u.role = sqlc.narg('role')::user_role)
  AND (sqlc.narg('status')::user_status IS NULL OR u.status = sqlc.narg('status')::user_status);

-- name: GetActiveUsersByRoles :many
SELECT * FROM users u
WHERE u.role = ANY(sqlc.arg('roles')::user_role[])
  AND u.status = 'active'
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, role, status) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUserRole :one
UPDATE users SET role = $2 WHERE id = $1 RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users SET email = $2 WHERE id = $1 RETURNING *;

-- name: UpdateUserFullName :one
UPDATE users SET full_name = $2 WHERE id = $1 RETURNING *;

-- name: UpdateUserStatus :one
UPDATE users SET status = $2 WHERE id = $1 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET password_hash = $2 WHERE id = $1 RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;
