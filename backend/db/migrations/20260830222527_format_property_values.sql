-- Create "format_property_value" function
CREATE FUNCTION "public"."format_property_value" ("property_value" jsonb, "value_type" text) RETURNS text LANGUAGE sql IMMUTABLE AS $$
SELECT CASE
    WHEN property_value IS NULL THEN ''
    WHEN value_type IN ('mass', 'volume') THEN
      trim_scale((property_value->>'amount')::NUMERIC / 10000)::TEXT
        || ' ' || (property_value->>'unit')
    WHEN value_type = 'price' THEN
      trim_scale((property_value->>'amount')::NUMERIC / 10000)::TEXT
        || ' ' || (property_value->>'currency')
    ELSE COALESCE(property_value #>> '{}', '')
  END;
$$;
-- Modify "render_item_derived_name" function
CREATE OR REPLACE FUNCTION "public"."render_item_derived_name" ("target_item_id" bigint, "derived_format" text) RETURNS text LANGUAGE sql STABLE AS $$
SELECT btrim(COALESCE((
    SELECT string_agg(
      CASE
        WHEN token.parts[1] IS NOT NULL
          THEN format_property_value(item_property.property_value, properties.value_type)
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
