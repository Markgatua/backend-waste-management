-- role_has_permissions.sql

-- name: GetRolePermissions :many
SELECT * FROM role_has_permissions WHERE role_id = $1;

-- name: AssignPermission :exec
insert into role_has_permissions (role_id,permission_id) VALUES($1,$2);

-- name: RevokePermission :exec
DELETE from role_has_permissions where role_id=$1 AND permission_id=$2;

-- name: GetDuplicateRoleHasPermission :one
select count(*) from role_has_permissions where role_id=$1 and permission_id=$2;