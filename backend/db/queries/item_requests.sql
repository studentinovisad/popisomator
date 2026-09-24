-- name: CreateItemRequest :one
INSERT INTO item_requests(user_id, item_id, status, reason)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: LockItemForRequest :one
SELECT id
FROM items
WHERE id = $1
FOR UPDATE;

-- name: HasApprovedItemRequest :one
SELECT EXISTS(
  SELECT 1
  FROM item_requests
  WHERE item_id = $1 AND status = 'approved'
);

-- name: CountItemRequests :one
SELECT count(*) FROM item_requests
WHERE (sqlc.narg('status')::request_status IS NULL OR item_requests.status = sqlc.narg('status'))
  AND (sqlc.narg('user_ids')::bigint[] IS NULL OR item_requests.user_id = ANY(sqlc.narg('user_ids')::bigint[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR item_requests.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR item_requests.created_at <= sqlc.narg('created_to'));

-- name: ListItemRequests :many
SELECT
  item_requests.*,
  users.full_name AS user_name,
  render_item_derived_name(items.id, item_types.derived_name_format) AS item_name
FROM item_requests
JOIN users ON users.id = item_requests.user_id
JOIN items ON items.id = item_requests.item_id
JOIN item_types ON item_types.id = items.type_id
WHERE (sqlc.narg('status')::request_status IS NULL OR item_requests.status = sqlc.narg('status'))
  AND (sqlc.narg('user_ids')::bigint[] IS NULL OR item_requests.user_id = ANY(sqlc.narg('user_ids')::bigint[]))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR item_requests.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR item_requests.created_at <= sqlc.narg('created_to'))
ORDER BY item_requests.created_at DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: ListItemRequestUsers :many
SELECT DISTINCT users.id, users.full_name
FROM item_requests
JOIN users ON users.id = item_requests.user_id
ORDER BY users.full_name, users.id;

-- name: ListItemPreparationRequests :many
SELECT
  item_requests.*,
  users.full_name AS user_name,
  render_item_derived_name(items.id, item_types.derived_name_format) AS item_name,
  item_types.name AS item_type_name,
  item_types.derived_name_format,
  items.consumption,
  items.location_id
FROM item_requests
JOIN users ON users.id = item_requests.user_id
JOIN items ON items.id = item_requests.item_id
JOIN item_types ON item_types.id = items.type_id
WHERE item_requests.user_id = $1
  AND (sqlc.narg('item_ids')::bigint[] IS NULL OR item_requests.item_id = ANY(sqlc.narg('item_ids')::bigint[]))
  AND (sqlc.narg('status')::request_status IS NULL OR item_requests.status = sqlc.narg('status'))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR item_requests.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR item_requests.created_at <= sqlc.narg('created_to'))
ORDER BY item_requests.created_at, item_requests.item_id;

-- name: GetItemRequest :one
SELECT sqlc.embed(item_requests),
  users.full_name AS user_name,
  render_item_derived_name(items.id, item_types.derived_name_format) AS item_name
FROM item_requests
JOIN users ON users.id = item_requests.user_id
JOIN items ON items.id = item_requests.item_id
JOIN item_types ON item_types.id = items.type_id
WHERE item_requests.user_id = $1 AND item_requests.item_id = $2;

-- name: CheckItemsForRequests :many
SELECT * FROM item_requests
WHERE item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: GetItemsRequestStatuses :many
SELECT 
  item_requests.item_id, 
  item_requests.status, 
  item_requests.user_id, 
  users.full_name AS user_full_name
FROM item_requests
JOIN users ON users.id = item_requests.user_id
WHERE (user_id = sqlc.arg('viewer_id') OR item_requests.status = 'approved')
  AND item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: UpdateItemRequest_Status :one
UPDATE item_requests
SET status = $3
WHERE user_id = $1 AND item_id = $2 AND status != $3
RETURNING *;

-- name: DeleteItemRequest :execrows
DELETE FROM item_requests WHERE user_id = $1 AND item_id = $2;

-- name: DeleteNonApprovedItemRequests :execrows
DELETE FROM item_requests WHERE item_id = $1 AND status = 'requested';
