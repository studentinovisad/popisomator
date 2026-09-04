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

-- name: ListNotifications :many
SELECT 
    sqlc.embed(notif), 
    sqlc.embed(notifdesc_item_request),
    sqlc.embed(notifdesc_item_expiry)
FROM notifications AS notif
LEFT JOIN notifdesc_item_request 
    ON notif.id = notifdesc_item_request.notification_id
LEFT JOIN notifdesc_item_expiry 
    ON notif.id = notifdesc_item_expiry.notification_id
WHERE recipient_id = $1
ORDER BY notif.created_at
LIMIT sqlc.arg('page_limit') OFFSET sqlc.arg('page_offset');

-- name: CountNotifications :one
SELECT count(*) FROM notifications
WHERE recipient_id = $1;

-- name: CountUnreadNotifications :one
SELECT count(*) FROM notifications
WHERE recipient_id = $1 AND read = true;

-- name: ReadNotifications :execrows
UPDATE notifications
SET read = true
WHERE read = false AND recipient_id = $1;

-- name: DeleteNotification :execrows
DELETE FROM notifications WHERE id = $1 AND recipient_id = $2;