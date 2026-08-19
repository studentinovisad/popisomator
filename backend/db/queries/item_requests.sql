-- name: CreateItemRequest :one
INSERT INTO item_requests(user_id, item_id, status, reason) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: CountItemRequests :one
SELECT count(*) FROM item_requests;

-- name: ListItemRequests :many
SELECT * FROM item_requests
WHERE (sqlc.narg('status')::request_status IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('user_id')::bigint IS NULL OR user_id = sqlc.narg('user_id'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: GetItemRequest :one
SELECT * FROM item_requests 
WHERE user_id = $1 AND item_id = $2;

-- name: CheckItemsForRequests :many
SELECT * FROM item_requests 
WHERE item_id = ANY(sqlc.arg('item_ids')::bigint[]);

-- name: ApproveItemRequest :one
UPDATE item_requests 
SET status = 'approved' 
WHERE user_id = $1 AND item_id = $2 AND status = 'requested' 
RETURNING *;

-- name: DeleteItemRequest :execrows
DELETE FROM item_requests WHERE user_id = $1 AND item_id = $2;

-- name: DeleteNonApprovedItemRequests :execrows
DELETE FROM item_requests WHERE item_id = $1 AND status = 'requested';