-- Modify "episodes" table
ALTER TABLE "public"."episodes"
  ALTER COLUMN "release_date" TYPE timestamptz USING CASE
    WHEN "release_date" IS NULL OR btrim("release_date") = '' THEN NULL
    ELSE "release_date"::date AT TIME ZONE 'UTC'
  END,
  ADD COLUMN "description" text NULL;
