-- regions.sql

-- name: GetAllOrganizations :many
SELECT * from organizations;

-- name: GetOrganization :one
SELECT * FROM organizations WHERE ID = $1;

-- name: InsertOrganization :one
insert into
    organizations(name, country_id)
values($1, $2) returning *;

-- name: GetOrganizationCountWithNameAndCountry :many
SELECT *
from organizations
where
    LOWER(name) = $1
    and country_id = $2;

-- name: GetDuplicateOrganization :many
SELECT *
FROM organizations
where
    id != $1
    and LOWER(name) = $2
    and country_id = $3;

-- name: UpdateOrganization :exec
update organizations set name=$1,country_id=$2 where id=$3;