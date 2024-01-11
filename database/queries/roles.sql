-- roles.sql

-- name: UpdateRole :exec
update roles set name=$2, guard_name=$3, description=$4, deleted_at=$5 where id=$1;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1 AND deleted_at IS NULL;

-- name: GetRoles :many
select * from roles where deleted_at IS NULL;

-- name: InsertRole :exec
insert into roles (name,guard_name,description) VALUES($1,$2,$3);
