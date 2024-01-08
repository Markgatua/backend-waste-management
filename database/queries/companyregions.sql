-- -- companiescompany_regionals.sql

-- -- name: GetAllCompaniesRegions :many
-- SELECT company_regionals.*, companies.name from company_regionals JOIN companies ON company_regionals.company_id = companies.id;

-- -- name: GetAllCompanyRegions :many
-- SELECT company_regionals.*, companies.name from company_regionals JOIN companies ON company_regionals.company_id = companies.id where company_regionals.company_id=$1;

-- -- name: GetOneCompanyRegion :one
-- SELECT company_regionals.*, companies.name FROM company_regionals JOIN companies ON company_regionals.company_id = companies.id WHERE company_regionals.id = $1;

-- -- name: InsertCompanyRegion :one
-- INSERT INTO company_regionals (company_id, region ) VALUES ($1, $2) RETURNING *;

-- -- name: UpdateCompanyRegionData :exec
-- update company_regionals set region=$2 where id=$1;
