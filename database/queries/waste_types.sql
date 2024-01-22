-- waste_types.sql

-- name: GetAllWasteTypes :many
select * from waste_types;

-- name: GetOneWasteType :one
select * from waste_types where id=$1;

-- name: InsertWasteType :one
INSERT INTO waste_types (name,category) VALUES ($1,$2) RETURNING *;

-- name: UpdateWasteType :exec
update waste_types set name=$2, category=$3, deleted_at=$4 where id=$1;

-- name: GetUsersWasteType :many
select * from waste_types where deleted_at is NULL;
