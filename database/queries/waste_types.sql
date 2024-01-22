-- waste_types.sql

-- name: GetAllWasteTypes :many
select waste_types.*,uploads.path as file_path from waste_types left join uploads on uploads.item_id=waste_types.id and uploads.related_table='waste_types';

-- name: GetOneWasteType :one
select * from waste_types where id=$1;

-- name: InsertWasteType :one
INSERT INTO waste_types (name,category) VALUES ($1,$2) RETURNING *;

-- name: UpdateWasteType :exec
update waste_types set name=$1, category=$2,is_active=$3 where id=$4;

-- name: GetUsersWasteType :many
select * from waste_types where deleted_at is NULL;
