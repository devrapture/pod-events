-- Modify "episodes" table
ALTER TABLE "public"."episodes"
  ALTER COLUMN "release_date" TYPE timestamptz USING CASE
    WHEN "release_date" IS NULL OR btrim("release_date") = '' THEN NULL
    WHEN btrim("release_date") ~ '^[0-9]{4}$' THEN (btrim("release_date") || '-01-01')::date AT TIME ZONE 'UTC'
    WHEN btrim("release_date") ~ '^[0-9]{4}-[0-9]{2}$' THEN (btrim("release_date") || '-01')::date AT TIME ZONE 'UTC'
    ELSE btrim("release_date")::date AT TIME ZONE 'UTC'
  END,
  ADD COLUMN "description" text NULL;
