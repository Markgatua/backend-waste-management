-- name: CreatePermission :exec
INSERT INTO
    permissions (
        permission_id,
        name,
        guard_name,
        module,
        submodule
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    );

-- name: GetAllPermissions :many
SELECT * FROM permissions;