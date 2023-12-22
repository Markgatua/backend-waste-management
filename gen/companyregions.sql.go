// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: companyregions.sql

package gen

import (
	"context"
	"time"
)

const getAllCompaniesRegions = `-- name: GetAllCompaniesRegions :many

SELECT company_regionals.id, company_regionals.company_id, company_regionals.region, company_regionals.created_at, companies.name from company_regionals JOIN companies ON company_regionals.company_id = companies.id
`

type GetAllCompaniesRegionsRow struct {
	ID        int32     `json:"id"`
	CompanyID int32     `json:"company_id"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

// companiescompany_regionals.sql
func (q *Queries) GetAllCompaniesRegions(ctx context.Context) ([]GetAllCompaniesRegionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllCompaniesRegions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllCompaniesRegionsRow{}
	for rows.Next() {
		var i GetAllCompaniesRegionsRow
		if err := rows.Scan(
			&i.ID,
			&i.CompanyID,
			&i.Region,
			&i.CreatedAt,
			&i.Name,
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

const getAllCompanyRegions = `-- name: GetAllCompanyRegions :many
SELECT company_regionals.id, company_regionals.company_id, company_regionals.region, company_regionals.created_at, companies.name from company_regionals JOIN companies ON company_regionals.company_id = companies.id where company_regionals.company_id=$1
`

type GetAllCompanyRegionsRow struct {
	ID        int32     `json:"id"`
	CompanyID int32     `json:"company_id"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

func (q *Queries) GetAllCompanyRegions(ctx context.Context, companyID int32) ([]GetAllCompanyRegionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllCompanyRegions, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllCompanyRegionsRow{}
	for rows.Next() {
		var i GetAllCompanyRegionsRow
		if err := rows.Scan(
			&i.ID,
			&i.CompanyID,
			&i.Region,
			&i.CreatedAt,
			&i.Name,
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

const getOneCompanyRegion = `-- name: GetOneCompanyRegion :one
SELECT company_regionals.id, company_regionals.company_id, company_regionals.region, company_regionals.created_at, companies.name FROM company_regionals JOIN companies ON company_regionals.company_id = companies.id WHERE company_regionals.id = $1
`

type GetOneCompanyRegionRow struct {
	ID        int32     `json:"id"`
	CompanyID int32     `json:"company_id"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

func (q *Queries) GetOneCompanyRegion(ctx context.Context, id int32) (GetOneCompanyRegionRow, error) {
	row := q.db.QueryRowContext(ctx, getOneCompanyRegion, id)
	var i GetOneCompanyRegionRow
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.Region,
		&i.CreatedAt,
		&i.Name,
	)
	return i, err
}

const insertCompanyRegion = `-- name: InsertCompanyRegion :one
INSERT INTO company_regionals (company_id, region ) VALUES ($1, $2) RETURNING id, company_id, region, created_at
`

type InsertCompanyRegionParams struct {
	CompanyID int32  `json:"company_id"`
	Region    string `json:"region"`
}

func (q *Queries) InsertCompanyRegion(ctx context.Context, arg InsertCompanyRegionParams) (CompanyRegional, error) {
	row := q.db.QueryRowContext(ctx, insertCompanyRegion, arg.CompanyID, arg.Region)
	var i CompanyRegional
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.Region,
		&i.CreatedAt,
	)
	return i, err
}

const updateCompanyRegionData = `-- name: UpdateCompanyRegionData :exec
update company_regionals set region=$2 where id=$1
`

type UpdateCompanyRegionDataParams struct {
	ID     int32  `json:"id"`
	Region string `json:"region"`
}

func (q *Queries) UpdateCompanyRegionData(ctx context.Context, arg UpdateCompanyRegionDataParams) error {
	_, err := q.db.ExecContext(ctx, updateCompanyRegionData, arg.ID, arg.Region)
	return err
}
