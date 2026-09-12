DROP TABLE IF EXISTS audit_log;
DROP TYPE IF EXISTS audit_target_type;
DROP TYPE IF EXISTS audit_action;

DROP FUNCTION IF EXISTS 
    render_item_derived_name(BIGINT, TEXT),
    escape_like_pattern(TEXT)
CASCADE;
DROP TABLE IF EXISTS 
    item_requests,
    item_properties,
    items,
    item_type_properties,
    item_types,
    properties,
    users,
    notifications,
    notifdesc_item_request,
    notifdesc_item_expiry
CASCADE;
DROP TYPE IF EXISTS 
    request_status,
    property_visibility,
    consumption_status,
    user_status,
    user_role,
    notification_kind,
    notifdesc_expiry_type
CASCADE;