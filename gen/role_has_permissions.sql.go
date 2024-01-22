// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: role_has_permissions.sql

package gen

import (
	"context"
)

const assignPermission = `-- name: AssignPermission :exec
insert into role_has_permissions (role_id,permission_id) VALUES($1,$2)
`

type AssignPermissionParams struct {
	RoleID       int32 `json:"role_id"`
	PermissionID int32 `json:"permission_id"`
}

func (q *Queries) AssignPermission(ctx context.Context, arg AssignPermissionParams) error {
	_, err := q.db.ExecContext(ctx, assignPermission, arg.RoleID, arg.PermissionID)
	return err
}

const getDuplicateRoleHasPermission = `-- name: GetDuplicateRoleHasPermission :one
select count(*) from role_has_permissions where role_id=$1 and permission_id=$2
`

type GetDuplicateRoleHasPermissionParams struct {
	RoleID       int32 `json:"role_id"`
	PermissionID int32 `json:"permission_id"`
}

func (q *Queries) GetDuplicateRoleHasPermission(ctx context.Context, arg GetDuplicateRoleHasPermissionParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDuplicateRoleHasPermission, arg.RoleID, arg.PermissionID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getRolePermissions = `-- name: GetRolePermissions :many

SELECT permission_id, role_id FROM role_has_permissions WHERE role_id = $1
`

// role_has_permissions.sql
func (q *Queries) GetRolePermissions(ctx context.Context, roleID int32) ([]RoleHasPermission, error) {
	rows, err := q.db.QueryContext(ctx, getRolePermissions, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RoleHasPermission{}
	for rows.Next() {
		var i RoleHasPermission
		if err := rows.Scan(&i.PermissionID, &i.RoleID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revokePermission = `-- name: RevokePermission :exec
DELETE from role_has_permissions where role_id=$1 AND permission_id=$2
`

type RevokePermissionParams struct {
	RoleID       int32 `json:"role_id"`
	PermissionID int32 `json:"permission_id"`
}

func (q *Queries) RevokePermission(ctx context.Context, arg RevokePermissionParams) error {
	_, err := q.db.ExecContext(ctx, revokePermission, arg.RoleID, arg.PermissionID)
	return err
}
