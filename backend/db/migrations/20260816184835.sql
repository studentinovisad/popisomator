-- Create enum type "property_visibility"
CREATE TYPE "public"."property_visibility" AS ENUM ('overview', 'details');
-- Modify "item_type_properties" table
ALTER TABLE "public"."item_type_properties" ADD COLUMN "visibility" "public"."property_visibility" NOT NULL DEFAULT 'overview';
-- Modify "item_types" table
ALTER TABLE "public"."item_types" ADD COLUMN "derived_name_format" text NULL;
