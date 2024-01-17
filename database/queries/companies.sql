-- companies.sql

-- name: GetAllCompanies :many
SELECT
  companies.*,
  organizations.name AS organization_name,
  counties.name AS county,
  sub_counties.name AS sub_county
FROM
  companies
LEFT JOIN
  organizations ON organizations.id = companies.organization_id
LEFT JOIN
  counties ON counties.id = companies.county_id
LEFT JOIN
  sub_counties ON sub_counties.id = companies.sub_county_id;


-- name: GetCompany :one
SELECT
  companies.*,
  organizations.name AS organization_name,
  counties.name AS county,
  sub_counties.name AS sub_county
FROM
  companies
LEFT JOIN
  organizations ON organizations.id = companies.organization_id
LEFT JOIN
  counties ON counties.id = companies.county_id
LEFT JOIN
  sub_counties ON sub_counties.id = companies.sub_county_id
WHERE companies.id = $1;

-- name: UpdateCompanyStatus :exec
update companies set is_active = $2 where id = $1;

-- name: InsertCompany :one
insert into
    companies(
        county_id,
        sub_county_id,
        physical_position,
        name,
        company_type,
        organization_id,
        region,
        location,
        is_active
    )
values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *;

-- name: UpdateCompany :exec
update companies
set
    county_id = $1,
    sub_county_id = $2,
    physical_position = $3,
    name = $4,
    company_type = $5,
    organization_id = $6,
    region = $7,
    location = $8,
    is_active = $9
where id = $10;

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

-- name: DeleteCompany :exec
delete from companies where id=$1;