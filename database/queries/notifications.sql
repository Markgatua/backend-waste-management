-- notifications.sql

-- name: InsertNewNotificationRequest :exec
insert into notifications( user_id,subject,message ) values ($1, $2, $3) returning *;

-- name: GetMyNotifications :many

SELECT * FROM notifications WHERE user_id = $1;

-- name: UpdateNotificationStatus :exec

update notifications set status=$2 where id=$1;
