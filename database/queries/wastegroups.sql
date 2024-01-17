-- waste_groups.sql

-- name: GetAllWasteGroups :many
select * from waste_groups;

-- name: GetOneWasteGroup :one
select * from waste_groups where id=$1;

-- name: InsertWasteGroup :one
INSERT INTO waste_groups (name,category) VALUES ($1,$2) RETURNING *;

-- name: UpdateWasteGroup :exec
update waste_groups set name=$2, category=$3, deleted_at=$4 where id=$1;

-- name: GetUsersWasteGroups :many
select * from waste_groups where deleted_at is NULL;
