-- Create "podcast_shows" table
CREATE TABLE "public"."podcast_shows" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "spotify_show_id" text NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "image_url" text NOT NULL,
  "spotify_url" text NOT NULL,
  "latest_episode_id" text NULL,
  "latest_episode_published_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_podcast_shows_deleted_at" to table: "podcast_shows"
CREATE INDEX "idx_podcast_shows_deleted_at" ON "public"."podcast_shows" ("deleted_at");
-- Create index "idx_podcast_shows_spotify_show_id" to table: "podcast_shows"
CREATE UNIQUE INDEX "idx_podcast_shows_spotify_show_id" ON "public"."podcast_shows" ("spotify_show_id");
-- Create "episodes" table
CREATE TABLE "public"."episodes" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "podcast_show_id" uuid NOT NULL,
  "spotify_episode_id" text NOT NULL,
  "audio_preview_url" text NULL,
  "spotify_url" text NULL,
  "duration_ms" bigint NULL,
  "release_date" text NULL,
  "image_url" text NULL,
  "name" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_podcast_shows_episode" FOREIGN KEY ("podcast_show_id") REFERENCES "public"."podcast_shows" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_episodes_deleted_at" to table: "episodes"
CREATE INDEX "idx_episodes_deleted_at" ON "public"."episodes" ("deleted_at");
-- Create index "idx_episodes_podcast_show_id" to table: "episodes"
CREATE INDEX "idx_episodes_podcast_show_id" ON "public"."episodes" ("podcast_show_id");
-- Create index "idx_episodes_spotify_episode_id" to table: "episodes"
CREATE UNIQUE INDEX "idx_episodes_spotify_episode_id" ON "public"."episodes" ("spotify_episode_id");
-- Create "notification_logs" table
CREATE TABLE "public"."notification_logs" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "episode_id" uuid NOT NULL,
  "channel_type" text NOT NULL,
  "status" text NOT NULL,
  "error_message" text NULL,
  "sent_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_logs_episode" FOREIGN KEY ("episode_id") REFERENCES "public"."episodes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notification_logs_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notification_logs_deleted_at" to table: "notification_logs"
CREATE INDEX "idx_notification_logs_deleted_at" ON "public"."notification_logs" ("deleted_at");
-- Create index "idx_notification_logs_episode_id" to table: "notification_logs"
CREATE INDEX "idx_notification_logs_episode_id" ON "public"."notification_logs" ("episode_id");
-- Create index "idx_notification_logs_user_id" to table: "notification_logs"
CREATE INDEX "idx_notification_logs_user_id" ON "public"."notification_logs" ("user_id");
-- Create "subscriptions" table
CREATE TABLE "public"."subscriptions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "podcast_show_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_podcast_shows_subscription" FOREIGN KEY ("podcast_show_id") REFERENCES "public"."podcast_shows" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_subscriptions" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subscriptions_deleted_at" to table: "subscriptions"
CREATE INDEX "idx_subscriptions_deleted_at" ON "public"."subscriptions" ("deleted_at");
-- Create index "idx_user_show" to table: "subscriptions"
CREATE UNIQUE INDEX "idx_user_show" ON "public"."subscriptions" ("user_id", "podcast_show_id");
