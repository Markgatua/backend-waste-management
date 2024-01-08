-- regions.sql

-- name: GetAllOrganizations :many
SELECT * from organizations;

-- name: GetOneOrganization :one
SELECT * FROM organizations WHERE ID = $1;

-- name: InsertOrganization :one
insert into organizations(name, country_id) values($1, $2) returning *;

-- name: GetOrganizationCountWithNameAndCountry :many
SELECT count(*) from organizations where LOWER(name) = LOWER($1) and country_id = $2;