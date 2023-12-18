-- name: GetAllUsers :many
select * from users;

-- name: GetUser :one
select * from users where id=$1;

-- name: UpdateUser :exec
update users set first_name=$1,last_name=$2 where id=$3;

-- name: GetUsersWithRole :many
select users.*, roles.name from users INNER JOIN roles ON users.role_id=roles.id;
