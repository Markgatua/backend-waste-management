-- name: GetAllMainOrganizationUsers :many
select users.id, users.first_name, users.last_name, users.email, users.avatar_url, users.calling_code, users.phone, users.is_active, roles.name as role_name,roles.id as role_id
from users
    inner join roles on users.role_id = roles.id
where
    users.email not ilike 'superadmin@admin.com' and users.is_main_organization_user=true;


-- name: GetAllUsers :many
select users.id, users.first_name, users.last_name, users.email, users.avatar_url, users.calling_code, users.phone, users.is_active, roles.name as role_name,roles.id as role_id
from users
    inner join roles on users.role_id = roles.id
where
    users.email not ilike 'superadmin@admin.com' and users.is_main_organization_user=false;

-- name: GetMainOrganizationUser :one
select * from users where id = $1;

-- name: UpdateUserIsActive :exec
update users set is_active=$1 where id =$2;

-- name: UpdateMainOrganizationUser :exec
update users set first_name = $1, last_name = $2 where id = $3;

-- name: GetMainOrganizationUserByEmail :one
select * from users where email = $1;

-- name: GetUsersWithRole :many
select users.*, roles.name from users INNER JOIN roles ON users.role_id = roles.id;

-- name: CreateMainOrganizationAdmin :exec
insert into users (first_name, last_name, email, provider, role_id, password, confirmed_at,is_main_organization_user) VALUES ($1, $2, $3, $4, $5, $6, $7,$8);