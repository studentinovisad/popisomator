CREATE TYPE user_role AS ENUM ('admin', 'manager', 'user');
CREATE TYPE user_status AS ENUM ('requested', 'active');

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    status user_status NOT NULL DEFAULT 'active'
);

CREATE TYPE consumption_status AS ENUM ('not_consumed', 'partially_consumed', 'fully_consumed', 'damaged');

CREATE TABLE properties (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    value_type VARCHAR(32) NOT NULL,
    default_value JSONB
);

CREATE TABLE item_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    derived_name_format TEXT,
    expiring_soon_days SMALLINT CHECK (expiring_soon_days > 0)
);

CREATE TYPE property_visibility AS ENUM ('overview', 'details');
CREATE TABLE item_type_properties (
    type_id BIGINT REFERENCES item_types(id) ON DELETE CASCADE,
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    default_value JSONB,
    visibility property_visibility NOT NULL DEFAULT 'overview',
    position INTEGER NOT NULL,
    PRIMARY KEY (type_id, property_id)
);

CREATE INDEX idx_item_type_properties_property_id 
  ON item_type_properties(property_id);
CREATE UNIQUE INDEX idx_item_type_properties_type_position 
  ON item_type_properties(type_id, position);

CREATE TABLE items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    consumption consumption_status NOT NULL DEFAULT 'not_consumed',
    type_id BIGINT NOT NULL REFERENCES item_types(id) ON DELETE RESTRICT
);

CREATE INDEX idx_items_type_id ON items(type_id);
CREATE INDEX idx_items_created_at ON items(created_at);

CREATE TABLE item_properties (
    item_id BIGINT REFERENCES items(id) ON DELETE CASCADE,
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    property_value JSONB NOT NULL,
    PRIMARY KEY (item_id, property_id)
);

CREATE INDEX idx_item_properties_property_id ON item_properties(property_id);

CREATE FUNCTION escape_like_pattern(value TEXT)
RETURNS TEXT
LANGUAGE sql
IMMUTABLE
AS $$
  SELECT replace(replace(replace(value, '\', '\\'), '%', '\%'), '_', '\_');
$$;

CREATE FUNCTION render_item_derived_name(target_item_id BIGINT, derived_format TEXT)
RETURNS TEXT
LANGUAGE sql
STABLE
AS $$
  SELECT btrim(COALESCE((
    SELECT string_agg(
      CASE
        WHEN token.parts[1] IS NOT NULL THEN COALESCE(item_property.property_value #>> '{}', '')
        ELSE token.parts[2]
      END,
      '' ORDER BY token.position
    )
    FROM regexp_matches(
      derived_format,
      '\{([^{}]+)\}|([^{}]+)',
      'g'
    ) WITH ORDINALITY AS token(parts, position)
    LEFT JOIN properties ON properties.name = btrim(token.parts[1])
    LEFT JOIN item_properties AS item_property
      ON item_property.item_id = target_item_id
     AND item_property.property_id = properties.id
  ), ''));
$$;

CREATE TYPE request_status AS ENUM ('requested', 'approved');
CREATE TABLE item_requests (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    item_id BIGINT REFERENCES items(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status request_status NOT NULL DEFAULT 'requested',
    reason TEXT NOT NULL,
    PRIMARY KEY (user_id, item_id)
);

CREATE INDEX idx_item_requests_user_id ON item_requests(user_id);
CREATE INDEX idx_item_requests_item_id ON item_requests(item_id);

CREATE UNIQUE INDEX idx_unique_approved_item_requests
ON item_requests(item_id)
WHERE status = 'approved';

CREATE TYPE notification_kind AS ENUM ('item_request', 'item_expiry');

CREATE TABLE notifications (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  kind notification_kind NOT NULL,
  read BOOLEAN NOT NULL DEFAULT false,
  UNIQUE (id, kind)
);

CREATE INDEX idx_notifications_recipient_id ON notifications(recipient_id);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

CREATE TABLE notifdesc_item_request (
  notification_id BIGINT PRIMARY KEY REFERENCES notifications(id) ON DELETE CASCADE,
  kind notification_kind GENERATED ALWAYS AS ('item_request') STORED,
  user_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,

  FOREIGN KEY(notification_id, kind)
    REFERENCES notifications(id, kind) ON DELETE CASCADE,
  FOREIGN KEY(user_id, item_id)
    REFERENCES item_requests(user_id, item_id) ON DELETE CASCADE
);

CREATE TYPE notifdesc_expiry_type AS ENUM ('expiring_soon', 'expired');

CREATE TABLE notifdesc_item_expiry (
  notification_id BIGINT PRIMARY KEY REFERENCES notifications(id) ON DELETE CASCADE,
  kind notification_kind GENERATED ALWAYS AS ('item_expiry') STORED,
  item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  expiry_type notifdesc_expiry_type NOT NULL,

  FOREIGN KEY(notification_id, kind)
    REFERENCES notifications(id, kind) ON DELETE CASCADE
);

-- An audit entry has to outlive the thing it describes: an item deleted a year ago is exactly the
-- history worth keeping. So target_id deliberately carries no foreign key, and target_label holds
-- the name the target went by at the time, which is impossible to recompute once it is gone and
-- misleading once it has been renamed. actor_id does keep a foreign key, because a user row going
-- away should not orphan the log - it nulls the link and leaves actor_name, the snapshot taken when
-- the change was made, as the only name the entry needs.

CREATE TYPE audit_action AS ENUM (
    'item_create', 'item_update', 'item_consume', 'item_delete',
    'item_property_add', 'item_property_update', 'item_property_remove',
    'item_request_create', 'item_request_approve', 'item_request_delete',
    'item_type_create', 'item_type_update', 'item_type_delete',
    'item_type_property_add', 'item_type_property_update', 'item_type_property_remove',
    'item_type_property_reorder',
    'property_create', 'property_update', 'property_delete'
);

-- The entity an entry is filed under. Item requests and the two property join tables are absent on
-- purpose: all three are keyed by a pair of columns, which a single target_id cannot hold. They are
-- filed under the item or item type they concern, with the other half of the key in context. That
-- is also what puts request activity on an item's own timeline, which is where the record of who
-- held an item has to live - approving a request deletes the competing ones, so item_requests
-- cannot answer that on its own.
CREATE TYPE audit_target_type AS ENUM ('item', 'item_type', 'property');

CREATE TABLE audit_log (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    actor_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    actor_name TEXT NOT NULL,
    action audit_action NOT NULL,
    target_type audit_target_type NOT NULL,
    target_id BIGINT NOT NULL,
    target_label TEXT NOT NULL,
    -- Field diffs, [{key, label, value_type, old, new}]. Empty for actions that are not edits.
    changes JSONB NOT NULL DEFAULT '[]'::jsonb,
    -- Everything the entry needs that is not a diff: the subject user of a request, the type an item
    -- was made from, how many items one bulk add produced, what a deleted item carried.
    context JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- The global feed. id breaks ties because a bulk add stamps every one of its rows with the same
-- now(), which would otherwise make paging through them unstable.
CREATE INDEX idx_audit_log_created_at ON audit_log (created_at DESC, id DESC);
-- The per-item and per-type timelines.
CREATE INDEX idx_audit_log_target ON audit_log (target_type, target_id, created_at DESC, id DESC);
-- The actor filter, and the actor option list behind it.
CREATE INDEX idx_audit_log_actor ON audit_log (actor_id, created_at DESC, id DESC);
