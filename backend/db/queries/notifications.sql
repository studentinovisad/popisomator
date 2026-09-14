-- name: CreateNotifications :many
INSERT INTO notifications (recipient_id, kind) 
VALUES (unnest(sqlc.arg('recipient_ids')::bigint[]), $1)
RETURNING *;

-- name: CreateNotificationDescriptors_ItemRequest :many
INSERT INTO notifdesc_item_request (notification_id, user_id, item_id) 
VALUES (unnest(sqlc.arg('notification_ids')::bigint[]), $1, $2)
RETURNING *;

-- name: CreateNotificationDescriptors_ItemExpiry :many
INSERT INTO notifdesc_item_expiry (notification_id, item_id, expiry_type)
VALUES (unnest(sqlc.arg('notification_ids')::bigint[]), $1, $2)
RETURNING *;

-- name: CreateNotificationDescriptors_LowStock :many
INSERT INTO notifdesc_low_stock (notification_id, type_id, type_label, group_name, threshold, observed)
VALUES (unnest(sqlc.arg('notification_ids')::bigint[]), $1, $2, $3, $4, $5)
RETURNING *;

-- name: ListNotifications :many
SELECT
    sqlc.embed(notif),
    notifdesc_item_request.user_id AS item_request_user_id,
    notifdesc_item_request.item_id AS item_request_item_id,
    notifdesc_item_expiry.item_id AS item_expiry_item_id,
    notifdesc_item_expiry.expiry_type AS item_expiry_type,
    notifdesc_low_stock.type_id AS low_stock_type_id,
    notifdesc_low_stock.type_label AS low_stock_type_label,
    notifdesc_low_stock.group_name AS low_stock_group_name,
    notifdesc_low_stock.threshold AS low_stock_threshold,
    notifdesc_low_stock.observed AS low_stock_observed
FROM notifications AS notif
LEFT JOIN notifdesc_item_request
    ON notif.id = notifdesc_item_request.notification_id
LEFT JOIN notifdesc_item_expiry
    ON notif.id = notifdesc_item_expiry.notification_id
LEFT JOIN notifdesc_low_stock
    ON notif.id = notifdesc_low_stock.notification_id
WHERE recipient_id = $1
-- Unread first, then newest first. The id tiebreaker keeps pagination stable: notifications are
-- bulk-inserted, so a whole batch shares one created_at.
ORDER BY notif.read ASC, notif.created_at DESC, notif.id DESC
LIMIT sqlc.arg('page_limit') OFFSET sqlc.arg('page_offset');

-- name: CountNotifications :one
SELECT count(*) FROM notifications
WHERE recipient_id = $1;

-- name: CountUnreadNotifications :one
SELECT count(*) FROM notifications
WHERE recipient_id = $1 AND read = false;

-- name: ReadNotifications :execrows
UPDATE notifications
SET read = true
WHERE read = false AND recipient_id = $1;

-- name: DeleteNotification :execrows
DELETE FROM notifications WHERE id = $1 AND recipient_id = $2;