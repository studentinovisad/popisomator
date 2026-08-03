env "local" {
  src = "file://db/schema.sql"
  dev = "docker://postgres/18/dev"
  migration {
    dir = "file://db/migrations"
  }
  url = getenv("POPISOMATOR_MIGRATION_DSN")
}
