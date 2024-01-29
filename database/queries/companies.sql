-- companies.sql

-- name: GetAllAggregators :many
SELECT
    companies.*,
    organizations.name AS organization_name,
    users.first_name,
    users.id as user_id,
    users.last_name,
    users.email,
    uploads.path as file_path,
    countries.name AS country_name

FROM
    companies
    left JOIN uploads on uploads.item_id = companies.id and uploads.related_table = 'companies'
    LEFT JOIN organizations ON organizations.id = companies.organization_id
    left join users on users.user_company_id = companies.id and users.is_company_super_admin = true
    LEFT JOIN countries ON countries.id = companies.country_id
WHERE
    companies.company_type = 2;

-- name: GetAllGreenChampions :many
SELECT
    companies.*,
    organizations.name AS organization_name,
     users.first_name,
    users.id as user_id,
    users.last_name,
    users.email,
    uploads.path as file_path,
    countries.name AS country_name
FROM
    companies
    left JOIN uploads on uploads.item_id = companies.id
    and uploads.related_table = 'companies'
    LEFT JOIN organizations ON organizations.id = companies.organization_id
        left join users on users.user_company_id = companies.id and users.is_company_super_admin = true
    LEFT JOIN countries ON countries.id = companies.country_id
WHERE
    companies.company_type = 1;

-- name: GetCompany :one
SELECT companies.*, organizations.name AS organization_name --, counties.name AS county
FROM
    companies
    LEFT JOIN organizations ON organizations.id = companies.organization_id
    --LEFT JOIN counties ON counties.id = companies.county_id
WHERE
    companies.id = $1;
   

-- name: UpdateCompanyStatus :exec
update companies set is_active = $2 where id = $1;

-- name: InsertCompany :one
insert into
    companies (
        country_id, region, name, administrative_level_1_location, company_type, organization_id, location, is_active, lat, lng
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    ) returning *;

-- name: UpdateCompany :exec
update companies
set
    country_id = $1,
    name = $2,
    company_type = $3,
    organization_id = $4,
    region = $5,
    location = $6,
    is_active = $7,
    lat = $8,
    lng = $9,
    administrative_level_1_location = $10
where
    id = $11;

-- name: GetDuplicateCompanies :many
select * from companies where lower(name) = $1 and organization_id = $2;

-- name: GetDuplicateCompaniesWithoutID :many
select *
from companies
where
    id != $1
    and lower(name) = $2
    and organization_id = $3;

-- name: DeleteCompany :exec
delete from companies where id = $1;