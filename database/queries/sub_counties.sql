-- sub_counties.sql

-- name: InsertSubcounties :exec
INSERT INTO sub_counties (name) VALUES($1,$2);

-- name: ViewSubCounties :many
SELECT * FROM sub_counties;
