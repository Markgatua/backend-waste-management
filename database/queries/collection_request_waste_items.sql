-- waste_items.sql

-- name: InsertCollectionRequestWasteItem :one
INSERT INTO collection_request_waste_items (collection_request_id,waste_type_id,weight) VALUES ($1,$2,$3) RETURNING *;

-- name: DeleteWasteItemsForCollectionRequest :exec
delete from collection_request_waste_items where collection_request_id=$1;