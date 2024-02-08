-- name: CreateVehicleTypes :one
INSERT INTO vehicle_types (id,name,max_vehicle_weight,max_vehicle_height,description) VALUES (
        sqlc.arg('id'),
        sqlc.arg('name'),
        sqlc.arg('max_vehicle_weight'),
        sqlc.arg('max_vehicle_height'),
        sqlc.arg('description')
    ) ON CONFLICT(id) do update set name=EXCLUDED.name,max_vehicle_weight=EXCLUDED.max_vehicle_weight,max_vehicle_height=EXCLUDED.max_vehicle_height,description=EXCLUDED.description returning *;

-- name: UpdateVehicle :exec
update vehicles set assigned_driver_id=$1,vehicle_type_id=$2,reg_no=$3,is_active=$4 where id=$5;

-- name: AddVehicle :one
insert into vehicles(company_id,assigned_driver_id,vehicle_type_id,reg_no,is_active) VALUES(
   sqlc.arg('company_id'),
   sqlc.arg('assigned_driver_id'),
   sqlc.arg('vehicle_type_id'),
   sqlc.arg('reg_no'),
   sqlc.arg('is_active')
)returning *;

-- name: GetAllVehicles :many
select vehicles.*,vehicle_types.description,vehicle_types.max_vehicle_height,vehicle_types.max_vehicle_weight,vehicle_types.name from vehicles inner join vehicle_types on vehicle_types.id=vehicles.vehicle_type_id where vehicles.company_id=$1;

-- name: DeleteVehicle :exec
delete from vehicles where id=$1;

-- name: GetDuplicateVehiclesWithoutID :one
select count(*) from vehicles where lower(reg_no)=$1 and id!=$2;

-- name: UpdateVehicleStatus :exec
update vehicles set is_active=$1 where id =$2;

-- name: GetDuplicateVehicle :one
select count(*) from vehicles where lower(reg_no)=$1;