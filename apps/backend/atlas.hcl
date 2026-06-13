data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "-tags=atlas",
    "./internal/migrations/loader",
  ]
}

env "local" {
  url = "postgres://postgres:password@127.0.0.1:5433/podevents?sslmode=disable"
  dev = "postgres://postgres:password@127.0.0.1:5433/podevents_dev?sslmode=disable"
  src = data.external_schema.gorm.url
  migration {
    dir = "file://migrations"
    revisions_schema = "public"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "production" {
  url = getenv("DATABASE_URL")
  dev = getenv("DEV_DATABASE_URL")
  src = data.external_schema.gorm.url
  migration {
    dir = "file://migrations"
    revisions_schema = "public"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
