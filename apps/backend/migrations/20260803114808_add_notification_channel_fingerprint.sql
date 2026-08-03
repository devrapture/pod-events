-- Drop index "idx_notification_channels_destination" from table: "notification_channels"
DROP INDEX "public"."idx_notification_channels_destination";
-- Modify "notification_channels" table
ALTER TABLE "public"."notification_channels" ADD COLUMN "destination_fingerprint" character(64) NULL;
-- Backfill existing rows with stable, unique placeholders before enforcing NOT NULL.
-- Webhook fingerprints can be replaced with keyed fingerprints by an application-level backfill.
UPDATE "public"."notification_channels"
SET "destination_fingerprint" = repeat(replace("id"::text, '-', ''), 2);
ALTER TABLE "public"."notification_channels" ALTER COLUMN "destination_fingerprint" SET NOT NULL;
-- Create index "idx_user_channel_destination" to table: "notification_channels"
CREATE UNIQUE INDEX "idx_user_channel_destination" ON "public"."notification_channels" ("user_id", "channel_type", "destination_fingerprint");
