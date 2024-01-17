-- sub_counties.sql

-- name: InsertSubcounties :exec
INSERT INTO sub_counties (name,county_id) VALUES($1,$2);

-- name: ViewSubCounties :many
SELECT * FROM sub_counties;


-- name: CheckSubCountiesDuplicate :one
SELECT COUNT(*) FROM sub_counties WHERE name=$1;


-- name: GetSubCountiesForACounty :many
SELECT * FROM sub_counties WHERE county_id=$1;