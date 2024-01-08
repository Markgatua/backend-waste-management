env "local" {
  src = "file://database/migrations/schema.sql"
  url = "postgres://mac:ServBay.dev@localhost/ttnm?sslmode=disable"
  dev = "postgres://mac:ServBay.dev@localhost/ttnm_dev?sslmode=disable"
  
  migration {
    dir = "file://database/migrations"
  }
}