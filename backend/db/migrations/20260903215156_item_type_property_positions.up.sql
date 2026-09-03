ALTER TABLE item_type_properties
ADD COLUMN position INTEGER;

WITH ordered_properties AS (
    SELECT
        type_id,
        property_id,
        row_number() OVER (PARTITION BY type_id ORDER BY property_id) - 1 AS position
    FROM item_type_properties
)
UPDATE item_type_properties
SET position = ordered_properties.position
FROM ordered_properties
WHERE item_type_properties.type_id = ordered_properties.type_id
  AND item_type_properties.property_id = ordered_properties.property_id;

ALTER TABLE item_type_properties
ALTER COLUMN position SET NOT NULL;

CREATE UNIQUE INDEX idx_item_type_properties_type_position
ON item_type_properties(type_id, position);
