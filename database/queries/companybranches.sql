-- company_branches.sql

-- name: GetAllBranches :many
SELECT company_branches.*, companies.name, company_regionals.region
FROM company_branches
JOIN companies ON companies.id = company_branches.company_id
LEFT JOIN company_regionals ON company_regionals.id = company_branches.region_id;

-- name: GetAllCompanyBranches :many
SELECT company_branches.*, companies.name, company_regionals.region
FROM company_branches
JOIN companies ON companies.id = company_branches.company_id
LEFT JOIN company_regionals ON company_regionals.id = company_branches.region_id
WHERE company_branches.company_id = $1;

-- name: GetBranchesForARegion :many
SELECT company_branches.*, companies.name, company_regionals.region
FROM company_branches
JOIN companies ON companies.id = company_branches.company_id
JOIN company_regionals ON company_regionals.id = company_branches.region_id
WHERE company_branches.region_id = $1;

-- name: GetOneCompanyBranch :one
SELECT company_branches.*, companies.name, company_regionals.region
FROM company_branches
JOIN companies ON companies.id = company_branches.company_id
JOIN company_regionals ON company_regionals.id = company_branches.region_id
WHERE company_branches.id = $1;

-- name: InsertCompanyBranch :one
INSERT INTO company_branches (company_id, region_id, branch, branch_location) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCompanyBranchStatus :exec
update company_branches set is_active=$2 where id=$1;

-- name: UpdateCompanyBranchData :exec
update company_branches set branch=$2,branch_location=$3,region_id=$4 where id=$1;
