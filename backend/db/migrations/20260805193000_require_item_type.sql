-- Items are always classified by one item type. Existing rows must be assigned a type before
-- this migration is applied.
ALTER TABLE "public"."items" ALTER COLUMN "type_id" SET NOT NULL;
ALTER TABLE "public"."items" DROP CONSTRAINT "items_type_id_fkey";
ALTER TABLE "public"."items" ADD CONSTRAINT "items_type_id_fkey"
  FOREIGN KEY ("type_id") REFERENCES "public"."item_types" ("id") ON DELETE RESTRICT;
