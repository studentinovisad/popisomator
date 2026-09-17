-- Stock is an aggregate, never a stored number: items holds one row per physical item and nothing
-- records a quantity. What makes two of those rows the same stock is their rendered derived name, so
-- that is what the count groups by.
--
-- The grouping deliberately spans every item of the type whatever its state, while only items that
-- are actually available count as stock. That split is the whole point. A group that has been used up
-- has no available rows left, and a plain GROUP BY over in-stock items would drop it from the result
-- entirely - losing precisely the group worth warning about. Counting inside a FILTER instead keeps
-- the consumed rows present as evidence the group exists, and reports it at zero.
--
-- Available means untouched and on the shelf: an item someone holds an approved request for is spoken
-- for and cannot be handed to anyone else, so it is not stock however full it still is. The join
-- mirrors the one ListItems filters by; idx_unique_approved_item_requests caps it at one row per
-- item, which is what keeps it from inflating total_count.
--
-- A null type_id counts the whole inventory instead of one type, which is the same question asked of
-- everything at once. The type rides along in the select because a group name only identifies a
-- stock line within one type, and the threshold because it is per type and there is otherwise no way
-- to tell a short group from a healthy one.
--
-- Deliberately unlimited. One caller wants the scarcest groups and another the largest, which are
-- opposite ends of this ordering, so a LIMIT would serve the first and quietly truncate the second.
--
-- name: ListStockGroups :many
SELECT
  item_types.id AS type_id,
  item_types.name AS type_name,
  item_types.low_stock_count,
  render_item_derived_name(items.id, item_types.derived_name_format) AS group_name,
  count(*) FILTER (
    WHERE items.consumption = 'not_consumed' AND approved_request.item_id IS NULL
  ) AS in_stock_count,
  count(*) AS total_count
FROM items
JOIN item_types ON item_types.id = items.type_id
LEFT JOIN item_requests AS approved_request
  ON approved_request.item_id = items.id
 AND approved_request.status = 'approved'
WHERE (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
GROUP BY item_types.id, group_name
ORDER BY in_stock_count, item_types.name, group_name;

-- name: GroupItemCountsForGroups :many
WITH requested_groups AS (
  SELECT DISTINCT requested.group_name
  FROM unnest(sqlc.arg('group_names')::text[]) AS requested(group_name)
)
SELECT
  requested_groups.group_name::text AS group_name,
  count(items.id) FILTER (
    WHERE items.consumption = 'not_consumed' AND approved_request.item_id IS NULL
  ) AS in_stock_count,
  count(items.id) AS total_count
FROM requested_groups
JOIN item_types ON item_types.id = sqlc.arg('type_id')
LEFT JOIN items
  ON items.type_id = item_types.id
 AND render_item_derived_name(items.id, item_types.derived_name_format) = requested_groups.group_name
LEFT JOIN item_requests AS approved_request
  ON approved_request.item_id = items.id
 AND approved_request.status = 'approved'
GROUP BY requested_groups.group_name;

-- name: ListLowStockAlertsForGroups :many
SELECT * FROM low_stock_alerts
WHERE type_id = $1 AND group_name = ANY(sqlc.arg('group_names')::text[]);

-- name: ListLowStockAlerts :many
SELECT * FROM low_stock_alerts
WHERE type_id = $1;

-- name: InsertLowStockAlert :execrows
INSERT INTO low_stock_alerts (type_id, group_name)
VALUES ($1, $2)
ON CONFLICT (type_id, group_name) DO NOTHING;

-- name: DeleteLowStockAlerts :execrows
DELETE FROM low_stock_alerts
WHERE type_id = $1 AND group_name = ANY(sqlc.arg('group_names')::text[]);

-- name: ClearLowStockAlerts :execrows
DELETE FROM low_stock_alerts WHERE type_id = $1;
