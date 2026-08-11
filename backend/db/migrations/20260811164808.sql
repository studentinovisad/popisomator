-- Create enum type "user_status"
CREATE TYPE "public"."user_status" AS ENUM ('requested', 'active');
-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "created_at", DROP COLUMN "updated_at", ADD COLUMN "status" "public"."user_status" NOT NULL DEFAULT 'active';
