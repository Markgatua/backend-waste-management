-- waste_items.sql

-- name: InsertWasteItem :one
INSERT INTO waste_items (collection_request_id,waste_group_id,weight) VALUES ($1,$2,$3) RETURNING *;
