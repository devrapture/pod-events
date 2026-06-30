-- Create index "idx_notification_channels_destination" to table: "notification_channels"
CREATE UNIQUE INDEX "idx_notification_channels_destination" ON "public"."notification_channels" ("destination");
