-- name: CreatePermission :exec
INSERT INTO
    permissions (
        name,
        guard_name,
        module,
        submodule
    )
VALUES (
        $1,
        $2,
        $3,
        $4
    );

-- name: GetAllPermissions :many
SELECT * FROM permissions;