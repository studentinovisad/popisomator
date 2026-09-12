-- The actor's display name is snapshotted here rather than looked up first, so recording a change
-- costs no extra round trip. An empty actor_name means there was no user behind the change at all -
-- the seeder, or a future background job - which the client words however it likes. That is distinct
-- from a null actor_id with a name, which is a user who has since been deleted.
-- name: WriteAuditEntry :exec
INSERT INTO audit_log (actor_id, actor_name, action, target_type, target_id, target_label, changes, context)
VALUES (
  sqlc.narg('actor_id')::bigint,
  COALESCE((SELECT full_name FROM users WHERE id = sqlc.narg('actor_id')::bigint), ''),
  sqlc.arg('action')::audit_action,
  sqlc.arg('target_type')::audit_target_type,
  sqlc.arg('target_id')::bigint,
  sqlc.arg('target_label')::text,
  sqlc.arg('changes')::jsonb,
  sqlc.arg('context')::jsonb
);

-- One action against several targets in a single statement, the way CreateNotifications inserts one
-- row per recipient. Used by CreateItem, which makes up to 100 items in a call and gives each its
-- own entry so that its timeline starts with its creation.
-- name: WriteAuditEntries :exec
INSERT INTO audit_log (actor_id, actor_name, action, target_type, target_id, target_label, changes, context)
SELECT
  sqlc.narg('actor_id')::bigint,
  COALESCE((SELECT full_name FROM users WHERE id = sqlc.narg('actor_id')::bigint), ''),
  sqlc.arg('action')::audit_action,
  sqlc.arg('target_type')::audit_target_type,
  entry.target_id,
  entry.target_label,
  entry.changes,
  entry.context
FROM ROWS FROM (
  unnest(sqlc.arg('target_ids')::bigint[]),
  unnest(sqlc.arg('target_labels')::text[]),
  unnest(sqlc.arg('changes')::jsonb[]),
  unnest(sqlc.arg('contexts')::jsonb[])
) AS entry(target_id, target_label, changes, context);

-- name: ListAuditLog :many
SELECT * FROM audit_log
WHERE (sqlc.narg('action')::audit_action IS NULL OR action = sqlc.narg('action'))
  AND (sqlc.narg('target_type')::audit_target_type IS NULL OR target_type = sqlc.narg('target_type'))
  AND (sqlc.narg('target_id')::bigint IS NULL OR target_id = sqlc.narg('target_id'))
  AND (sqlc.narg('actor_id')::bigint IS NULL OR actor_id = sqlc.narg('actor_id'))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to'))
-- Newest first. The id tiebreaker keeps pagination stable: a bulk add writes a whole batch of rows
-- under one created_at.
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountAuditLog :one
SELECT count(*) FROM audit_log
WHERE (sqlc.narg('action')::audit_action IS NULL OR action = sqlc.narg('action'))
  AND (sqlc.narg('target_type')::audit_target_type IS NULL OR target_type = sqlc.narg('target_type'))
  AND (sqlc.narg('target_id')::bigint IS NULL OR target_id = sqlc.narg('target_id'))
  AND (sqlc.narg('actor_id')::bigint IS NULL OR actor_id = sqlc.narg('actor_id'))
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to'));

-- Everyone who has ever made a recorded change, for the actor filter. Mirrors ListItemRequestUsers.
--
-- Driven from users rather than from audit_log: DISTINCT ON cannot skip-scan in Postgres, so
-- reading it the other way round sorts the whole log - a table that only ever grows and is never
-- pruned - to produce a handful of rows. EXISTS turns that into one index lookup per user against
-- idx_audit_log_actor, so the cost tracks the number of users instead of the size of the history.
--
-- Deleted users drop out either way: their actor_id is nulled, so there is no id left to filter by.
-- Their entries stay in the log and still show the name they acted under - they just cannot be
-- singled out by this filter any more.
-- name: ListAuditLogActors :many
SELECT users.id, users.full_name AS name
FROM users
WHERE EXISTS (SELECT 1 FROM audit_log WHERE audit_log.actor_id = users.id)
ORDER BY users.full_name, users.id;
