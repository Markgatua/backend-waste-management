-- companies.sql

-- name: GetAllCompanies :many
select * from companies;

-- name: GetCompany :one
select * from companies where id=$1;

-- name: InsertCompany :one
INSERT INTO companies (name, companytype, logo, location, has_regions, has_branches) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateCompanyStatus :exec
update companies set is_active=$2 where id=$1;

-- name: UpdateCompanyData :exec
update companies set name=$2,location=$3,logo=$4, has_regions=$5, has_branches=$6 where id=$1;
