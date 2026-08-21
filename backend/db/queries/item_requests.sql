-- name: CreateItemRequest :one
INSERT INTO item_requests(user_id, item_id, status, reason) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: CountItemRequests :one
SELECT count(*) FROM item_requests
WHERE (sqlc.narg('status')::request_status IS NULL OR item_requests.status = sqlc.narg('status'))
  AND (sqlc.narg('user_id')::bigint IS NULL OR item_requests.user_id = sqlc.narg('user_id'));

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
  AND (sqlc.narg('user_id')::bigint IS NULL OR item_requests.user_id = sqlc.narg('user_id'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: GetItemRequest :one
SELECT * FROM item_requests 
WHERE user_id = $1 AND item_id = $2;

-- name: CheckItemsForRequests :many
SELECT * FROM item_requests 
WHERE item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: GetUserItemRequests :many
SELECT * FROM item_requests
WHERE user_id = sqlc.arg('user_id')
  AND item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: ApproveItemRequest :one
UPDATE item_requests 
SET status = 'approved' 
WHERE user_id = $1 AND item_id = $2 AND status = 'requested' 
RETURNING *;

-- name: DeleteItemRequest :execrows
DELETE FROM item_requests WHERE user_id = $1 AND item_id = $2;

-- name: DeleteNonApprovedItemRequests :execrows
DELETE FROM item_requests WHERE item_id = $1 AND status = 'requested';
