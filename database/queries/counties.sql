-- counties.sql

-- name: InsertCounties :exec
INSERT INTO counties (name) VALUES($1);

-- name: ViewCounties :many
SELECT * FROM counties;

-- name: DuplicateCounties :one
SELECT COUNT(*) FROM counties WHERE name=$1;