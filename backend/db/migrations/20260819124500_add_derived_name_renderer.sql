CREATE FUNCTION "public"."render_item_derived_name"(
  "target_item_id" bigint,
  "derived_format" text
) RETURNS text
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
