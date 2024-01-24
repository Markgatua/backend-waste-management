-- waste_types.sql

-- name: GetAllWasteTypes :many
select waste_types.*,uploads.path as file_path from waste_types left join uploads on uploads.item_id=waste_types.id and uploads.related_table='waste_types';

-- name: GetMainWasteTypes :many
select waste_types.*,uploads.path as file_path from waste_types left join uploads on uploads.item_id=waste_types.id and uploads.related_table='waste_types' where waste_types.parent_id is null;

-- name: GetChildrenWasteTypes :many
select waste_types.*,uploads.path as file_path from waste_types left join uploads on uploads.item_id=waste_types.id and uploads.related_table='waste_types' where waste_types.parent_id =$1;


-- name: GetOneWasteType :one
select * from waste_types where id=$1;

-- name: InsertWasteType :one
INSERT INTO waste_types (name,category,parent_id) VALUES ($1,$2,$3) RETURNING *;

-- name: UpdateWasteType :exec
update waste_types set name=$1, category=$2,is_active=$3,parent_id=$4 where id=$5;

-- name: GetUsersWasteType :many
select * from waste_types where deleted_at is NULL;
