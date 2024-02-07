-- pickup_time_stamps.sql

-- name: InsertPickupTimeStsmp :exec
insert into pickup_time_stamps( id, stamp,time_range ) values (
        sqlc.arg('id'),
        sqlc.arg('stamp'),
        sqlc.arg('time_range')
    ) ON CONFLICT(id) do update set stamp=EXCLUDED.stamp,time_range=EXCLUDED.time_range returning *;


-- name: GetPickupTimeStamps :many

select * from pickup_time_stamps;

-- name: UpdatePickupTimeStamp :exec
update pickup_time_stamps
set
    stamp = $1,
    time_range = $2
where
    id = $3;

-- name: DuplicatePickupTimeStamp :one
SELECT COUNT(*) FROM pickup_time_stamps WHERE stamp=$1;