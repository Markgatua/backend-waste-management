-- regions.sql

-- name: GetAllOrganizations :many
SELECT organizations.*,countries.name as country from organizations left join countries on countries.id=organizations.country_id;

-- name: GetOrganization :one
SELECT organizations.*,countries.name as country FROM organizations left join countries on countries.id=organizations.country_id WHERE organizations.id = $1;

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
update organizations set name = $1, country_id = $2 where id = $3;

-- name: DeleteOrganization :exec
delete from organizations where id = $1;