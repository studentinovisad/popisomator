-- Create "escape_like_pattern" function
CREATE FUNCTION "public"."escape_like_pattern" ("value" text) RETURNS text LANGUAGE sql IMMUTABLE AS $$ SELECT replace(replace(replace(value, '\', '\\'), '%', '\%'), '_', '\_'); $$;
