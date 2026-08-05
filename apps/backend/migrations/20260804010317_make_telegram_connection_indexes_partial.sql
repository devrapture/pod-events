-- Drop index "idx_telegram_connections_token_hash" from table: "telegram_connections"
DROP INDEX "public"."idx_telegram_connections_token_hash";
-- Drop index "idx_telegram_connections_user_id" from table: "telegram_connections"
DROP INDEX "public"."idx_telegram_connections_user_id";
-- Create index "idx_telegram_connections_token_hash" to table: "telegram_connections"
CREATE UNIQUE INDEX "idx_telegram_connections_token_hash" ON "public"."telegram_connections" ("token_hash") WHERE (deleted_at IS NULL);
-- Create index "idx_telegram_connections_user_id" to table: "telegram_connections"
CREATE UNIQUE INDEX "idx_telegram_connections_user_id" ON "public"."telegram_connections" ("user_id") WHERE (deleted_at IS NULL);
