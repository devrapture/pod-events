-- Create "spotify_tokens" table
CREATE TABLE "public"."spotify_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "access_token" text NOT NULL,
  "refresh_token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "scope" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_spotify_token" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_spotify_tokens_deleted_at" to table: "spotify_tokens"
CREATE INDEX "idx_spotify_tokens_deleted_at" ON "public"."spotify_tokens" ("deleted_at");
-- Create index "idx_spotify_tokens_user_id" to table: "spotify_tokens"
CREATE UNIQUE INDEX "idx_spotify_tokens_user_id" ON "public"."spotify_tokens" ("user_id");
