-- An audit entry has to outlive the thing it describes: an item deleted a year ago is exactly the
-- history worth keeping. So target_id deliberately carries no foreign key, and target_label holds
-- the name the target went by at the time, which is impossible to recompute once it is gone and
-- misleading once it has been renamed. actor_id does keep a foreign key, because a user row going
-- away should not orphan the log - it nulls the link and leaves actor_name, the snapshot taken when
-- the change was made, as the only name the entry needs.

CREATE TYPE audit_action AS ENUM (
    'item_create', 'item_update', 'item_consume', 'item_delete',
    'item_property_add', 'item_property_update', 'item_property_remove',
    'item_request_create', 'item_request_approve', 'item_request_delete', 'item_request_supersede',
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
