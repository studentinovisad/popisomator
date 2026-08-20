CREATE TYPE user_role AS ENUM ('admin', 'manager', 'user');
CREATE TYPE user_status AS ENUM ('requested', 'active');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
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
    derived_name_format TEXT
);

CREATE TYPE property_visibility AS ENUM ('overview', 'details');
CREATE TABLE item_type_properties (
    type_id BIGINT REFERENCES item_types(id) ON DELETE CASCADE,
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    default_value JSONB,
    visibility property_visibility NOT NULL DEFAULT 'overview',
    PRIMARY KEY (type_id, property_id)
);

CREATE INDEX idx_item_type_properties_property_id ON item_type_properties(property_id);

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
