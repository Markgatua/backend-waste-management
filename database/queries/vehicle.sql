-- name: CreateVehicleTypes :one
INSERT INTO vehicle_types (id,name,max_vehicle_weight,max_vehicle_height,description) VALUES (
        sqlc.arg('id'),
        sqlc.arg('name'),
        sqlc.arg('max_vehicle_weight'),
        sqlc.arg('max_vehicle_height'),
        sqlc.arg('description')
    ) ON CONFLICT(id) do update set name=EXCLUDED.name,max_vehicle_weight=EXCLUDED.max_vehicle_weight,max_vehicle_height=EXCLUDED.max_vehicle_height,description=EXCLUDED.description returning *;
