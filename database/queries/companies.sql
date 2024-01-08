-- companies.sql

-- name: GetAllCompanies :many
select * from companies;

-- name: GetCompany :one
select * from companies where id = $1;

-- name: UpdateCompanyStatus :exec
update companies set is_active = $2 where id = $1;

-- name: InsertCompany :one
insert into
    companies(
        name,
        company_type,
        organization_id,
        region,
        location,
        is_active
    )
values ($1, $2, $3, $4, $5, $6) returning *;

-- name: UpdateCompany :exec
update companies
set
    name = $1,
    company_type = $2,
    organization_id = $3,
    region = $4,
    location = $5,
    is_active = $6
where id = $7;

-- name: GetDuplicateCompanies :many
select *
from companies
where
    lower(name) = $1
    and organization_id = $2;

-- name: GetDuplicateCompaniesWithoutID :many
select *
from companies
where
    id = $1
    and lower(name) = $2
    and organization_id = $3;