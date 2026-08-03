-- Modify "episodes" table
ALTER TABLE "public"."episodes" DROP COLUMN "audio_preview_url";
-- Drop index "idx_user_channel_destination" from table: "notification_channels"
DROP INDEX "public"."idx_user_channel_destination";
-- Create index "idx_user_channel_destination" to table: "notification_channels"
CREATE UNIQUE INDEX "idx_user_channel_destination" ON "public"."notification_channels" ("user_id", "channel_type", "destination_fingerprint") WHERE (deleted_at IS NULL);
