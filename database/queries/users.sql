-- name: GetAllUsers :many
select * from users;

-- name: GetUser :one
select * from users where id=$1;

-- name: UpdateUser :exec
update users set first_name=$1,last_name=$2 where id=$3;

-- name: GetUserByEmail :one
select * from users where email = $1;
-- name: GetUsersWithRole :many
select users.*, roles.name from users INNER JOIN roles ON users.role_id=roles.id;

-- name: CreateAdmin :exec
insert into users (first_name,last_name,email,provider,role_id,password,confirmed_at) VALUES($1,$2,$3,$4,$5,$6,$7);