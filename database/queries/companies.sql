-- companies.sql

-- name: GetAllAggregators :many
SELECT
  companies.*,
  organizations.name AS organization_name,
  uploads.path as file_path,
  counties.name AS county
FROM companies
left JOIN uploads on uploads.item_id=companies.id and uploads.related_table='companies'
LEFT JOIN organizations ON organizations.id = companies.organization_id
LEFT JOIN counties ON counties.id = companies.county_id
WHERE companies.company_type=2;

-- name: GetAllGreenChampions :many
SELECT
  companies.*,
  organizations.name AS organization_name,
  uploads.path as file_path,
  counties.name AS county
FROM companies
left JOIN uploads on uploads.item_id=companies.id and uploads.related_table='companies'
LEFT JOIN organizations ON organizations.id = companies.organization_id
LEFT JOIN counties ON counties.id = companies.county_id
WHERE companies.company_type=2;

-- name: GetCompany :one
SELECT
  companies.*,
  organizations.name AS organization_name,
  counties.name AS county
FROM
  companies
LEFT JOIN
  organizations ON organizations.id = companies.organization_id
LEFT JOIN
  counties ON counties.id = companies.county_id
WHERE companies.id = $1;

-- name: UpdateCompanyStatus :exec
update companies set is_active = $2 where id = $1;

-- name: InsertCompany :one
insert into
    companies(
        county_id,
        physical_position,
        name,
        company_type,
        organization_id,
        region,
        location,
        is_active
    )
values ($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: UpdateCompany :exec
update companies
set
    county_id = $1,
    physical_position = $2,
    name = $3,
    company_type = $4,
    organization_id = $5,
    region = $6,
    location = $7,
    is_active = $8
where id = $9;

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