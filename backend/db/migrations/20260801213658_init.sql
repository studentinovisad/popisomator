-- Create enum type "user_role"
CREATE TYPE "public"."user_role" AS ENUM ('admin', 'manager', 'user');
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "email" text NOT NULL,
  "password_hash" text NOT NULL,
  "full_name" text NOT NULL,
  "role" "public"."user_role" NOT NULL DEFAULT 'user',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
