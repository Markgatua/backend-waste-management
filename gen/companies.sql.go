// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: companies.sql

package gen

import (
	"context"
	"database/sql"
	"time"
)

const deleteCompany = `-- name: DeleteCompany :exec
delete from companies where id = $1
`

func (q *Queries) DeleteCompany(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCompany, id)
	return err
}

const getAllAggregators = `-- name: GetAllAggregators :many

SELECT
    companies.id, companies.name, companies.country_id, companies.company_type, companies.organization_id, companies.region, companies.location, companies.administrative_level_1_location, companies.lat, companies.lng, companies.is_active, companies.created_at,
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
    companies.company_type = 2
`

type GetAllAggregatorsRow struct {
	ID                           int32           `json:"id"`
	Name                         string          `json:"name"`
	CountryID                    int32           `json:"country_id"`
	CompanyType                  int32           `json:"company_type"`
	OrganizationID               sql.NullInt32   `json:"organization_id"`
	Region                       sql.NullString  `json:"region"`
	Location                     sql.NullString  `json:"location"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	IsActive                     bool            `json:"is_active"`
	CreatedAt                    time.Time       `json:"created_at"`
	OrganizationName             sql.NullString  `json:"organization_name"`
	FirstName                    sql.NullString  `json:"first_name"`
	UserID                       sql.NullInt32   `json:"user_id"`
	LastName                     sql.NullString  `json:"last_name"`
	Email                        sql.NullString  `json:"email"`
	FilePath                     sql.NullString  `json:"file_path"`
	CountryName                  sql.NullString  `json:"country_name"`
}

// companies.sql
func (q *Queries) GetAllAggregators(ctx context.Context) ([]GetAllAggregatorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAggregators)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllAggregatorsRow{}
	for rows.Next() {
		var i GetAllAggregatorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CountryID,
			&i.CompanyType,
			&i.OrganizationID,
			&i.Region,
			&i.Location,
			&i.AdministrativeLevel1Location,
			&i.Lat,
			&i.Lng,
			&i.IsActive,
			&i.CreatedAt,
			&i.OrganizationName,
			&i.FirstName,
			&i.UserID,
			&i.LastName,
			&i.Email,
			&i.FilePath,
			&i.CountryName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllGreenChampions = `-- name: GetAllGreenChampions :many
SELECT
    companies.id, companies.name, companies.country_id, companies.company_type, companies.organization_id, companies.region, companies.location, companies.administrative_level_1_location, companies.lat, companies.lng, companies.is_active, companies.created_at,
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
    companies.company_type = 1
`

type GetAllGreenChampionsRow struct {
	ID                           int32           `json:"id"`
	Name                         string          `json:"name"`
	CountryID                    int32           `json:"country_id"`
	CompanyType                  int32           `json:"company_type"`
	OrganizationID               sql.NullInt32   `json:"organization_id"`
	Region                       sql.NullString  `json:"region"`
	Location                     sql.NullString  `json:"location"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	IsActive                     bool            `json:"is_active"`
	CreatedAt                    time.Time       `json:"created_at"`
	OrganizationName             sql.NullString  `json:"organization_name"`
	FirstName                    sql.NullString  `json:"first_name"`
	UserID                       sql.NullInt32   `json:"user_id"`
	LastName                     sql.NullString  `json:"last_name"`
	Email                        sql.NullString  `json:"email"`
	FilePath                     sql.NullString  `json:"file_path"`
	CountryName                  sql.NullString  `json:"country_name"`
}

func (q *Queries) GetAllGreenChampions(ctx context.Context) ([]GetAllGreenChampionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllGreenChampions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllGreenChampionsRow{}
	for rows.Next() {
		var i GetAllGreenChampionsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CountryID,
			&i.CompanyType,
			&i.OrganizationID,
			&i.Region,
			&i.Location,
			&i.AdministrativeLevel1Location,
			&i.Lat,
			&i.Lng,
			&i.IsActive,
			&i.CreatedAt,
			&i.OrganizationName,
			&i.FirstName,
			&i.UserID,
			&i.LastName,
			&i.Email,
			&i.FilePath,
			&i.CountryName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompany = `-- name: GetCompany :one
SELECT companies.id, companies.name, companies.country_id, companies.company_type, companies.organization_id, companies.region, companies.location, companies.administrative_level_1_location, companies.lat, companies.lng, companies.is_active, companies.created_at, organizations.name AS organization_name, counties.name AS county
FROM
    companies
    LEFT JOIN organizations ON organizations.id = companies.organization_id
    LEFT JOIN counties ON counties.id = companies.county_id
WHERE
    companies.id = $1
`

type GetCompanyRow struct {
	ID                           int32           `json:"id"`
	Name                         string          `json:"name"`
	CountryID                    int32           `json:"country_id"`
	CompanyType                  int32           `json:"company_type"`
	OrganizationID               sql.NullInt32   `json:"organization_id"`
	Region                       sql.NullString  `json:"region"`
	Location                     sql.NullString  `json:"location"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	IsActive                     bool            `json:"is_active"`
	CreatedAt                    time.Time       `json:"created_at"`
	OrganizationName             sql.NullString  `json:"organization_name"`
	County                       sql.NullString  `json:"county"`
}

func (q *Queries) GetCompany(ctx context.Context, id int32) (GetCompanyRow, error) {
	row := q.db.QueryRowContext(ctx, getCompany, id)
	var i GetCompanyRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CountryID,
		&i.CompanyType,
		&i.OrganizationID,
		&i.Region,
		&i.Location,
		&i.AdministrativeLevel1Location,
		&i.Lat,
		&i.Lng,
		&i.IsActive,
		&i.CreatedAt,
		&i.OrganizationName,
		&i.County,
	)
	return i, err
}

const getDuplicateCompanies = `-- name: GetDuplicateCompanies :many
select id, name, country_id, company_type, organization_id, region, location, administrative_level_1_location, lat, lng, is_active, created_at from companies where lower(name) = $1 and organization_id = $2
`

type GetDuplicateCompaniesParams struct {
	Name           string        `json:"name"`
	OrganizationID sql.NullInt32 `json:"organization_id"`
}

func (q *Queries) GetDuplicateCompanies(ctx context.Context, arg GetDuplicateCompaniesParams) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, getDuplicateCompanies, arg.Name, arg.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Company{}
	for rows.Next() {
		var i Company
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CountryID,
			&i.CompanyType,
			&i.OrganizationID,
			&i.Region,
			&i.Location,
			&i.AdministrativeLevel1Location,
			&i.Lat,
			&i.Lng,
			&i.IsActive,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDuplicateCompaniesWithoutID = `-- name: GetDuplicateCompaniesWithoutID :many
select id, name, country_id, company_type, organization_id, region, location, administrative_level_1_location, lat, lng, is_active, created_at
from companies
where
    id != $1
    and lower(name) = $2
    and organization_id = $3
`

type GetDuplicateCompaniesWithoutIDParams struct {
	ID             int32         `json:"id"`
	Name           string        `json:"name"`
	OrganizationID sql.NullInt32 `json:"organization_id"`
}

func (q *Queries) GetDuplicateCompaniesWithoutID(ctx context.Context, arg GetDuplicateCompaniesWithoutIDParams) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, getDuplicateCompaniesWithoutID, arg.ID, arg.Name, arg.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Company{}
	for rows.Next() {
		var i Company
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CountryID,
			&i.CompanyType,
			&i.OrganizationID,
			&i.Region,
			&i.Location,
			&i.AdministrativeLevel1Location,
			&i.Lat,
			&i.Lng,
			&i.IsActive,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCompany = `-- name: InsertCompany :one
insert into
    companies (
        country_id, region, name, administrative_level_1_location, company_type, organization_id, location, is_active, lat, lng
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    ) returning id, name, country_id, company_type, organization_id, region, location, administrative_level_1_location, lat, lng, is_active, created_at
`

type InsertCompanyParams struct {
	CountryID                    int32           `json:"country_id"`
	Region                       sql.NullString  `json:"region"`
	Name                         string          `json:"name"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	CompanyType                  int32           `json:"company_type"`
	OrganizationID               sql.NullInt32   `json:"organization_id"`
	Location                     sql.NullString  `json:"location"`
	IsActive                     bool            `json:"is_active"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
}

func (q *Queries) InsertCompany(ctx context.Context, arg InsertCompanyParams) (Company, error) {
	row := q.db.QueryRowContext(ctx, insertCompany,
		arg.CountryID,
		arg.Region,
		arg.Name,
		arg.AdministrativeLevel1Location,
		arg.CompanyType,
		arg.OrganizationID,
		arg.Location,
		arg.IsActive,
		arg.Lat,
		arg.Lng,
	)
	var i Company
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CountryID,
		&i.CompanyType,
		&i.OrganizationID,
		&i.Region,
		&i.Location,
		&i.AdministrativeLevel1Location,
		&i.Lat,
		&i.Lng,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const updateCompany = `-- name: UpdateCompany :exec
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
    id = $11
`

type UpdateCompanyParams struct {
	CountryID                    int32           `json:"country_id"`
	Name                         string          `json:"name"`
	CompanyType                  int32           `json:"company_type"`
	OrganizationID               sql.NullInt32   `json:"organization_id"`
	Region                       sql.NullString  `json:"region"`
	Location                     sql.NullString  `json:"location"`
	IsActive                     bool            `json:"is_active"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	ID                           int32           `json:"id"`
}

func (q *Queries) UpdateCompany(ctx context.Context, arg UpdateCompanyParams) error {
	_, err := q.db.ExecContext(ctx, updateCompany,
		arg.CountryID,
		arg.Name,
		arg.CompanyType,
		arg.OrganizationID,
		arg.Region,
		arg.Location,
		arg.IsActive,
		arg.Lat,
		arg.Lng,
		arg.AdministrativeLevel1Location,
		arg.ID,
	)
	return err
}

const updateCompanyStatus = `-- name: UpdateCompanyStatus :exec
update companies set is_active = $2 where id = $1
`

type UpdateCompanyStatusParams struct {
	ID       int32 `json:"id"`
	IsActive bool  `json:"is_active"`
}

func (q *Queries) UpdateCompanyStatus(ctx context.Context, arg UpdateCompanyStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateCompanyStatus, arg.ID, arg.IsActive)
	return err
}
