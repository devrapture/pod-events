-- Allow existing Spotify-authenticated users to be linked to Google by email,
-- while new users can exist without an individual Spotify account.
ALTER TABLE "public"."users"
  ALTER COLUMN "spotify_user_id" DROP NOT NULL,
  ADD COLUMN "google_user_id" text NULL;

CREATE UNIQUE INDEX "idx_users_google_user_id" ON "public"."users" ("google_user_id");
