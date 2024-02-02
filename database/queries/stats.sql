
-- name: GetOrganizationCount :many
SELECT
    organizations.organization_type,
    SUM(organizations.id) AS total
FROM
    organizations
GROUP BY
    organizations.organization_type,organizations.id;



-- name: GetBranchCount :many
SELECT
    companies.company_type,
    SUM(companies.id) AS total
FROM
    companies
GROUP BY
    companies.company_type,companies.id;
