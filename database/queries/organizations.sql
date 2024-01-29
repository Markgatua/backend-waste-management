-- regions.sql

-- name: GetAllOrganizations :many
SELECT organizations.*,uploads.path as file_path,
users.first_name,users.id as user_id,users.last_name,users.email,
countries.id as country_id , countries.name as country from organizations 
left join countries on countries.id=organizations.country_id
left join users on users.user_organization_id = organizations.id and users.is_organization_super_admin=true
left join uploads on uploads.item_id=organizations.id and uploads.related_table='organizations';


-- name: UpdateOrganizationIsActive :exec
update organizations set is_active=$1 where id =$2;


-- name: GetOrganization :one
SELECT organizations.*,countries.name as country FROM organizations left join countries on countries.id=organizations.country_id WHERE organizations.id = $1;

-- name: InsertOrganization :one
insert into
    organizations(name, country_id,organization_type)
values($1, $2,$3) returning *;

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
update organizations set name = $1, country_id = $2, organization_type=$3 where id = $4;

-- name: DeleteOrganization :exec
delete from organizations where id = $1;