-- name: CreateCountry :exec
INSERT INTO
    countries (
        name,
        currency_code,
        capital,
        citizenship,
        country_code,
        currency,
        currency_sub_unit,
        currency_symbol,
        currency_decimals,
        full_name,
        iso_3166_2,
        iso_3166_3,
        region_code,
        sub_region_code,
        eea,
        calling_code,
        flag
    )
VALUES
(
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14,
        $15,
        $16,
        $17
    );

-- name: GetCountryBeCountryCode :many
SELECT * FROM countries WHERE country_code = $1;

-- name: GetCountryByName :one
select * from countries where name ilike sqlc.arg('country');

-- name: GetAllCountries :many
select * from countries;
