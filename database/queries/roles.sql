-- roles.sql

-- name: UpdateRole :exec
update roles set name=$2, guard_name=$3, description=$4, deleted_at=$5 where role_id=$1;

-- name: GetRole :one
SELECT * FROM roles WHERE role_id = $1 AND deleted_at IS NULL;

-- name: GetRoles :many
select * from roles where deleted_at IS NULL;

-- name: InsertRole :exec
insert into roles (name,role_id,guard_name,description) VALUES($1,$2,$3,$4);

-- name: GetDuplicateRole :one
select count(*) from roles where role_id=$1 and deleted_at IS null;