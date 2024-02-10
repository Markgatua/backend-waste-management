-- pickup_time_stamps.sql

-- name: InsertPickupTimeStsmp :exec
insert into pickup_time_stamps( id, stamp,time_range ,position) values (
        sqlc.arg('id'),
        sqlc.arg('stamp'),
        sqlc.arg('time_range'),
        sqlc.arg('position')
    ) ON CONFLICT(id) do update set stamp=EXCLUDED.stamp,time_range=EXCLUDED.time_range,position=EXCLUDED.position  returning *;


-- name: GetPickupTimeStamps :many

select * from pickup_time_stamps;

-- name: UpdatePickupTimeStamp :exec
update pickup_time_stamps
set
    stamp = $1,
    time_range = $2,
    position=$3
where
    id = $4;

-- name: DuplicatePickupTimeStamp :one
SELECT COUNT(*) FROM pickup_time_stamps WHERE stamp=$1;