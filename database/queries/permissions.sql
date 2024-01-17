-- name: CreatePermission :one
INSERT INTO
    permissions (
        name,
        action,
        module,
        submodule
    )
VALUES (
        sqlc.arg('name'),
        sqlc.arg('action'),
        sqlc.arg('module'),
        sqlc.arg('submodule')
    ) RETURNING *;

-- name: DeletePermissionByIds :exec
delete from permissions where not (id = ANY(sqlc.arg('permission_ids')::int[]));

-- name: DeletePermissionByActions :exec
delete from permissions where not (action = ANY(sqlc.arg('actions')::varchar[]));

-- name: GetAllPermissions :many
SELECT * FROM permissions;

-- name: GetPermissionsForRoleID :many
select
    permissions.id as permission_id,
    permissions.name,
    permissions.name,
    permissions.action
from role_has_permissions
    inner join permissions on permissions.id = role_has_permissions.permission_id
    where role_has_permissions.role_id=sqlc.arg('role_id');