env "local" {
  src = "file://database/migrations/schema.sql"
  url = "postgres://gakobo:Psql4321@localhost/ttnm_waste?sslmode=disable"
  dev = "postgres://gakobo:Psql4321@localhost/ttnm_waste_dev?sslmode=disable"

  migration {
    dir = "file://database/migrations"
  }
}