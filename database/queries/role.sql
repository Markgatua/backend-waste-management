-- name: CreateRole :one
INSERT INTO roles (id,name,guard_name,description,is_active) VALUES (
        sqlc.arg('id'),
        sqlc.arg('name'),
        sqlc.arg('guard_name'),
        sqlc.arg('description'),
        sqlc.arg('is_active')
    ) ON CONFLICT(id) do update set name=EXCLUDED.name,guard_name=EXCLUDED.guard_name,description=EXCLUDED.description,is_active=EXCLUDED.is_active returning *;



-- name: RoleExists :one
SELECT count(*) FROM roles where name = sqlc.arg('name');

-- name: GetRoles :many
select * from roles where id !=12;

-- name: GetRole :one
select * from roles where id=sqlc.arg('id');

-- name: DeleteRole :exec
delete from roles where id = sqlc.arg('id');

-- name: RemovePermissionsFromRole :exec
delete from role_has_permissions where permission_id = ANY(sqlc.arg('permission_ids'):: int []) and role_id=sqlc.arg('role_id');

-- name: DeactivateRole :exec
update roles set is_active = false where id = sqlc.arg('role_id');
-- name: ActivateRole :exec
update roles set is_active = true where id = sqlc.arg('role_id');

-- name: UpdateRole :exec
update roles set name = sqlc.arg('name'), is_active =sqlc.arg('is_active'),description = sqlc.arg('description') where id = sqlc.arg('role_id');

-- name: AssignPermissionToRole :exec
insert into role_has_permissions(role_id, permission_id) VALUES (sqlc.arg('role_id'),sqlc.arg('permission_id'));