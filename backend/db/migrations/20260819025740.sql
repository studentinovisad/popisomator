-- Create enum type "request_status"
CREATE TYPE "public"."request_status" AS ENUM ('requested', 'approved');
-- Create "item_requests" table
CREATE TABLE "public"."item_requests" (
  "user_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "status" "public"."request_status" NOT NULL DEFAULT 'requested',
  "reason" text NOT NULL,
  PRIMARY KEY ("user_id", "item_id"),
  CONSTRAINT "item_requests_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "public"."items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "item_requests_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_item_requests_item_id" to table: "item_requests"
CREATE INDEX "idx_item_requests_item_id" ON "public"."item_requests" ("item_id");
-- Create index "idx_item_requests_user_id" to table: "item_requests"
CREATE INDEX "idx_item_requests_user_id" ON "public"."item_requests" ("user_id");
-- Create index "idx_unique_approved_item_requests" to table: "item_requests"
CREATE UNIQUE INDEX "idx_unique_approved_item_requests" ON "public"."item_requests" ("item_id") WHERE (status = 'approved'::public.request_status);
