-- regions.sql

-- name: GetAllOrganizations :many
SELECT * from organizations;

-- name: GetOneOrganization :one

SELECT * FROM organizations WHERE ID=$1;
