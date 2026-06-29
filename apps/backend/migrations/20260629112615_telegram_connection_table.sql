-- Create "telegram_connections" table
CREATE TABLE "public"."telegram_connections" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "consumed" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_telegram_connections_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_telegram_connections_deleted_at" to table: "telegram_connections"
CREATE INDEX "idx_telegram_connections_deleted_at" ON "public"."telegram_connections" ("deleted_at");
-- Create index "idx_telegram_connections_token_hash" to table: "telegram_connections"
CREATE UNIQUE INDEX "idx_telegram_connections_token_hash" ON "public"."telegram_connections" ("token_hash");
-- Create index "idx_telegram_connections_user_id" to table: "telegram_connections"
CREATE UNIQUE INDEX "idx_telegram_connections_user_id" ON "public"."telegram_connections" ("user_id");
