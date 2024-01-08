-- companies.sql

-- name: GetAllCompanies :many
select * from companies;

-- name: GetCompany :one
select * from companies where id=$1;

-- name: UpdateCompanyStatus :exec
update companies set is_active=$2 where id=$1;

