-- name: GetAllUsers :many
select * from users;

-- name: GetUser :one
select * from users where id=$1;

-- name: UpdateUser :exec
update users set first_name=$1,last_name=$2 where id=$3;

-- name: GetUsersWithRole :many
select users.*, roles.name from users INNER JOIN roles ON users.role_id=roles.id;

-- name: CreateAdmin :exec
insert into users (first_name,last_name,email,provider,role_id,email,password) VALUES($1,$2,$3,$4,$5,$6,$7);

-- name: GetCompanyUsers :many
select users.*,companies.name as company,companies.location as companyLocation FROM users
LEFT JOIN
  companies ON companies.id = users.user_company_id
 where user_company_id = $1;