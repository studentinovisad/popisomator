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
  items.consumption
FROM item_requests
JOIN users ON users.id = item_requests.user_id
JOIN items ON items.id = item_requests.item_id
JOIN item_types ON item_types.id = items.type_id
WHERE item_requests.user_id = $1
  AND item_requests.status = 'requested'
ORDER BY item_requests.created_at, item_requests.item_id;

-- name: GetItemRequest :one
SELECT * FROM item_requests
WHERE user_id = $1 AND item_id = $2;

-- Every request standing against one item, with the requester's name. Used when an item is deleted:
-- the rows are about to cascade away, and each person who loses their claim gets their own audit
-- entry, so the name has to come back with them.
-- name: ListItemRequestsForItem :many
SELECT item_requests.user_id, item_requests.reason, item_requests.status,
       users.full_name AS user_name
FROM item_requests
JOIN users ON users.id = item_requests.user_id
WHERE item_requests.item_id = $1;

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

-- name: ApproveItemRequest :one
UPDATE item_requests
SET status = 'approved'
WHERE user_id = $1 AND item_id = $2 AND status = 'requested'
RETURNING *;

-- name: DeleteItemRequest :execrows
DELETE FROM item_requests WHERE user_id = $1 AND item_id = $2;

-- Approving a request cancels every other pending one for the item. The rows come back rather than
-- just their count, so the audit log can name each person whose request was superseded. The join to
-- users is safe: user_id is half the primary key and carries a foreign key, so it is never null and
-- never dangling.
-- name: DeleteNonApprovedItemRequests :many
WITH deleted AS (
  DELETE FROM item_requests
  WHERE item_id = $1 AND status = 'requested'
  RETURNING user_id, reason
)
SELECT deleted.user_id, deleted.reason, users.full_name AS user_name
FROM deleted
JOIN users ON users.id = deleted.user_id;
