-- An expiry is stored as a date string, and nothing stops a bad one getting in. The CASE is what
-- guards the cast - putting the check in the WHERE would not, since Postgres is free to evaluate
-- the two in either order.
--
-- The months are generated rather than read off the rows, so a period where nothing expires still
-- gets a bar. Generating them in the database keeps one clock in charge of which month is current.
-- name: CountExpiringByMonth :many
WITH axis AS (
  SELECT generate_series(
    date_trunc('month', current_date::timestamp),
    date_trunc('month', current_date::timestamp)
      + make_interval(months => sqlc.arg('months')::int - 1),
    interval '1 month'
  )::date AS month
),
expiring AS (
  SELECT date_trunc('month', expiry.expires_on::timestamp)::date AS month
  FROM items
  JOIN item_properties ON item_properties.item_id = items.id
  JOIN properties
    ON properties.id = item_properties.property_id
   AND properties.value_type = 'expiry'
  CROSS JOIN LATERAL (
    SELECT CASE
      WHEN pg_input_is_valid(item_properties.property_value #>> '{}', 'date')
      THEN (item_properties.property_value #>> '{}')::date
    END
  ) AS expiry(expires_on)
  WHERE items.consumption = 'not_consumed'
    AND expiry.expires_on IS NOT NULL
    AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
)
SELECT
  axis.month,
  count(expiring.month) AS item_count
FROM axis
LEFT JOIN expiring ON expiring.month = axis.month
GROUP BY axis.month
ORDER BY axis.month;

-- What a manager has to act on rather than only count: everything already past its date, plus what
-- falls inside the warning window its type sets. A type without a window warns about nothing, but an
-- item that is already expired is expired either way, which is why the two are separate conditions.
--
-- days_remaining comes back from here rather than being worked out from the date by the caller, so
-- the sign that decides expired from merely close is settled against one clock.
--
-- The list is capped, so every row carries how many there were before the cap. Counting over the
-- window rather than in a second query is what stops the total from ever disagreeing with the
-- condition above: there is only one copy of it. A cap that hid how much it was hiding would report
-- the ceiling as though it were the answer.
-- name: ListExpiringItems :many
WITH expiring AS (
  SELECT
    items.id,
    render_item_derived_name(items.id, item_types.derived_name_format) AS group_name,
    item_types.name AS type_name,
    expiry.expires_on::date AS expires_on,
    (expiry.expires_on - current_date)::int AS days_remaining
  FROM items
  JOIN item_types ON item_types.id = items.type_id
  JOIN item_properties ON item_properties.item_id = items.id
  JOIN properties
    ON properties.id = item_properties.property_id
   AND properties.value_type = 'expiry'
  CROSS JOIN LATERAL (
    SELECT CASE
      WHEN pg_input_is_valid(item_properties.property_value #>> '{}', 'date')
      THEN (item_properties.property_value #>> '{}')::date
    END
  ) AS expiry(expires_on)
  WHERE items.consumption = 'not_consumed'
    AND expiry.expires_on IS NOT NULL
    AND (
      expiry.expires_on < current_date
      OR (
        item_types.expiring_soon_days IS NOT NULL
        AND expiry.expires_on < current_date + make_interval(days => item_types.expiring_soon_days::int)
      )
    )
    AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
)
SELECT id, group_name, type_name, expires_on, days_remaining, count(*) OVER () AS total_count
FROM expiring
ORDER BY expires_on, group_name, id
LIMIT sqlc.arg('limit_val');

-- name: CountExpiredBacklog :one
SELECT count(*) FROM items
JOIN item_properties ON item_properties.item_id = items.id
JOIN properties
  ON properties.id = item_properties.property_id
 AND properties.value_type = 'expiry'
CROSS JOIN LATERAL (
  SELECT CASE
    WHEN pg_input_is_valid(item_properties.property_value #>> '{}', 'date')
    THEN (item_properties.property_value #>> '{}')::date
  END
) AS expiry(expires_on)
WHERE items.consumption = 'not_consumed'
  AND expiry.expires_on IS NOT NULL
  AND expiry.expires_on < current_date
  AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'));

-- An item records the state it is in, never when it got there, so the log is the only thing that
-- can date a consumption.
--
-- Consumption is a status rather than an event, and correcting one later writes a second entry.
-- Requiring the change to start at not_consumed counts each item the once it left the shelf. The
-- diff is stored JSON-encoded, which is why ->> hands back a bare enum name, and a consumption
-- entry only ever holds the one change.
--
-- The created_at bound repeats what the join already does, to keep the index in play.
--
-- Narrowing by type has to reach the item, and target_id has no foreign key to reach it with: an
-- item consumed and since deleted counts when no type is asked for, and drops out when one is.
-- name: CountConsumptionByMonth :many
WITH axis AS (
  SELECT generate_series(
    date_trunc('month', now()) - make_interval(months => sqlc.arg('months')::int - 1),
    date_trunc('month', now()),
    interval '1 month'
  )::date AS month
),
consumed AS (
  SELECT
    date_trunc('month', audit_log.created_at)::date AS month,
    audit_log.changes->0->>'new' AS state
  FROM audit_log
  LEFT JOIN items ON items.id = audit_log.target_id
  WHERE audit_log.action = 'item_consume'
    AND audit_log.target_type = 'item'
    AND audit_log.changes->0->>'old' = 'not_consumed'
    AND audit_log.changes->0->>'new' <> 'not_consumed'
    AND audit_log.created_at >= date_trunc('month', now())
      - make_interval(months => sqlc.arg('months')::int - 1)
    AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
)
SELECT
  axis.month,
  count(*) FILTER (WHERE consumed.state = 'fully_consumed') AS fully_consumed,
  count(*) FILTER (WHERE consumed.state = 'partially_consumed') AS partially_consumed,
  count(*) FILTER (WHERE consumed.state = 'damaged') AS damaged
FROM axis
LEFT JOIN consumed ON consumed.month = axis.month
GROUP BY axis.month
ORDER BY axis.month;

-- name: SumConsumptionQuantityByMonth :many
WITH axis AS (
  SELECT generate_series(
    date_trunc('month', now()) - make_interval(months => sqlc.arg('months')::int - 1),
    date_trunc('month', now()),
    interval '1 month'
  )::date AS month
),
consumed AS (
  SELECT items.id, items.type_id, date_trunc('month', audit_log.created_at)::date AS month
  FROM audit_log
  JOIN items ON items.id = audit_log.target_id
  WHERE audit_log.action = 'item_consume'
    AND audit_log.target_type = 'item'
    AND audit_log.changes->0->>'old' = 'not_consumed'
    AND audit_log.changes->0->>'new' <> 'not_consumed'
    AND audit_log.created_at >= date_trunc('month', now())
      - make_interval(months => sqlc.arg('months')::int - 1)
    AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
),
relevant_properties AS (
  SELECT DISTINCT item_type_properties.type_id, properties.id AS property_id, properties.name AS property_name, properties.value_type
  FROM item_type_properties
  JOIN properties ON properties.id = item_type_properties.property_id
  WHERE properties.value_type IN ('price', 'mass', 'volume')
    AND (sqlc.narg('type_id')::bigint IS NULL OR item_type_properties.type_id = sqlc.narg('type_id'))
),
sums AS (
  SELECT
    consumed.type_id,
    consumed.month,
    properties.id AS property_id,
    COALESCE(item_property.property_value ->> 'currency', '')::text AS currency,
    trim_scale(sum(
      (item_property.property_value ->> 'amount')::numeric * COALESCE(unit_factor.factor, 1)
    ))::text AS total_amount,
    count(*) AS value_count
  FROM consumed
  JOIN item_properties AS item_property ON item_property.item_id = consumed.id
  JOIN properties ON properties.id = item_property.property_id AND properties.value_type IN ('price', 'mass', 'volume')
  LEFT JOIN ROWS FROM (
    unnest(sqlc.arg('unit_value_types')::text[]),
    unnest(sqlc.arg('unit_names')::text[]),
    unnest(sqlc.arg('unit_factors')::bigint[])
  ) AS unit_factor(value_type, unit_name, factor)
    ON unit_factor.value_type = properties.value_type
   AND unit_factor.unit_name = item_property.property_value ->> 'unit'
  WHERE properties.value_type = 'price' OR unit_factor.factor IS NOT NULL
  GROUP BY consumed.type_id, consumed.month, properties.id, currency
)
SELECT
  relevant_properties.type_id,
  item_types.name AS type_name,
  axis.month,
  relevant_properties.property_id,
  relevant_properties.property_name,
  relevant_properties.value_type,
  COALESCE(sums.currency, '')::text AS currency,
  COALESCE(sums.total_amount, '0')::text AS total_amount,
  COALESCE(sums.value_count, 0)::bigint AS value_count
FROM relevant_properties
JOIN item_types ON item_types.id = relevant_properties.type_id
CROSS JOIN axis
LEFT JOIN sums
  ON sums.type_id = relevant_properties.type_id
 AND sums.month = axis.month
 AND sums.property_id = relevant_properties.property_id
ORDER BY item_types.name, relevant_properties.property_id, axis.month;

-- name: ListMostConsumedGroups :many
WITH consumed AS (
  SELECT
    items.id,
    items.type_id,
    item_types.name AS type_name,
    render_item_derived_name(items.id, item_types.derived_name_format) AS group_name
  FROM audit_log
  JOIN items ON items.id = audit_log.target_id
  JOIN item_types ON item_types.id = items.type_id
  WHERE audit_log.action = 'item_consume'
    AND audit_log.target_type = 'item'
    AND audit_log.changes->0->>'old' = 'not_consumed'
    AND audit_log.changes->0->>'new' <> 'not_consumed'
    AND audit_log.created_at >= date_trunc('month', now())
      - make_interval(months => sqlc.arg('months')::int - 1)
    AND (sqlc.narg('type_id')::bigint IS NULL OR items.type_id = sqlc.narg('type_id'))
),
ranked AS (
  SELECT
    type_id, type_name, group_name,
    count(*) AS consumed_count,
    row_number() OVER (ORDER BY count(*) DESC, group_name) AS rn
  FROM consumed
  GROUP BY type_id, type_name, group_name
),
top_groups AS (
  SELECT type_id, type_name, group_name, consumed_count
  FROM ranked
  WHERE rn <= sqlc.arg('limit_val')::int
),
sums AS (
  SELECT
    consumed.type_id,
    consumed.group_name,
    properties.id AS property_id,
    properties.name AS property_name,
    properties.value_type,
    COALESCE(item_property.property_value ->> 'currency', '')::text AS currency,
    trim_scale(sum(
      (item_property.property_value ->> 'amount')::numeric * COALESCE(unit_factor.factor, 1)
    ))::text AS total_amount,
    count(*) AS value_count
  FROM consumed
  JOIN top_groups ON top_groups.type_id = consumed.type_id AND top_groups.group_name = consumed.group_name
  JOIN item_properties AS item_property ON item_property.item_id = consumed.id
  JOIN properties ON properties.id = item_property.property_id AND properties.value_type IN ('price', 'mass', 'volume')
  LEFT JOIN ROWS FROM (
    unnest(sqlc.arg('unit_value_types')::text[]),
    unnest(sqlc.arg('unit_names')::text[]),
    unnest(sqlc.arg('unit_factors')::bigint[])
  ) AS unit_factor(value_type, unit_name, factor)
    ON unit_factor.value_type = properties.value_type
   AND unit_factor.unit_name = item_property.property_value ->> 'unit'
  WHERE properties.value_type = 'price' OR unit_factor.factor IS NOT NULL
  GROUP BY consumed.type_id, consumed.group_name, properties.id, properties.value_type, currency
)
SELECT
  top_groups.type_id,
  top_groups.type_name,
  top_groups.group_name,
  top_groups.consumed_count,
  sums.property_id,
  sums.property_name,
  sums.value_type,
  COALESCE(sums.currency, '')::text AS currency,
  sums.total_amount,
  sums.value_count
FROM top_groups
LEFT JOIN sums ON sums.type_id = top_groups.type_id AND sums.group_name = top_groups.group_name
ORDER BY top_groups.consumed_count DESC, top_groups.group_name, sums.property_id;
