-- Create "notification_channels" table
CREATE TABLE "public"."notification_channels" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "channel_type" text NOT NULL,
  "destination" text NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "label" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_notification_channels" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notification_channels_deleted_at" to table: "notification_channels"
CREATE INDEX "idx_notification_channels_deleted_at" ON "public"."notification_channels" ("deleted_at");
-- Create index "idx_notification_channels_user_id" to table: "notification_channels"
CREATE INDEX "idx_notification_channels_user_id" ON "public"."notification_channels" ("user_id");
