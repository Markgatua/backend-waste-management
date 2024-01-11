-- role_has_permissions.sql


-- name: GetRolePermissions :many
SELECT * FROM role_has_permissions WHERE role_id = $1;
