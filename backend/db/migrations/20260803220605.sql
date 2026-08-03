-- Create enum type "consumption_status"
CREATE TYPE "public"."consumption_status" AS ENUM ('not_consumed', 'partially_consumed', 'fully_consumed', 'damaged');
-- Create "item_types" table
CREATE TABLE "public"."item_types" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "name" text NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "item_types_name_key" UNIQUE ("name")
);
-- Create "items" table
CREATE TABLE "public"."items" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "consumption" "public"."consumption_status" NOT NULL DEFAULT 'not_consumed',
  "type_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "items_type_id_fkey" FOREIGN KEY ("type_id") REFERENCES "public"."item_types" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "properties" table
CREATE TABLE "public"."properties" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "name" text NOT NULL,
  "description" text NULL,
  "value_type" character varying(32) NOT NULL,
  "default_value" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "properties_name_key" UNIQUE ("name")
);
-- Create "item_properties" table
CREATE TABLE "public"."item_properties" (
  "item_id" bigint NOT NULL,
  "property_id" bigint NOT NULL,
  "property_value" jsonb NOT NULL,
  PRIMARY KEY ("item_id", "property_id"),
  CONSTRAINT "item_properties_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "public"."items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "item_properties_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "public"."properties" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_item_properties_property_id" to table: "item_properties"
CREATE INDEX "idx_item_properties_property_id" ON "public"."item_properties" ("property_id");
-- Create "item_type_properties" table
CREATE TABLE "public"."item_type_properties" (
  "type_id" bigint NOT NULL,
  "property_id" bigint NOT NULL,
  "default_value" jsonb NULL,
  PRIMARY KEY ("type_id", "property_id"),
  CONSTRAINT "item_type_properties_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "public"."properties" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "item_type_properties_type_id_fkey" FOREIGN KEY ("type_id") REFERENCES "public"."item_types" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_item_type_properties_property_id" to table: "item_type_properties"
CREATE INDEX "idx_item_type_properties_property_id" ON "public"."item_type_properties" ("property_id");
