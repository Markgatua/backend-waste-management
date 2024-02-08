// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: vehicle.sql

package gen

import (
	"context"
	"database/sql"
)

const addVehicle = `-- name: AddVehicle :one
insert into vehicles(company_id,assigned_driver_id,vehicle_type_id,reg_no,is_active) VALUES(
   $1,
   $2,
   $3,
   $4,
   $5
)returning id, company_id, assigned_driver_id, vehicle_type_id, reg_no, is_active
`

type AddVehicleParams struct {
	CompanyID        int32         `json:"company_id"`
	AssignedDriverID sql.NullInt32 `json:"assigned_driver_id"`
	VehicleTypeID    int32         `json:"vehicle_type_id"`
	RegNo            string        `json:"reg_no"`
	IsActive         bool          `json:"is_active"`
}

func (q *Queries) AddVehicle(ctx context.Context, arg AddVehicleParams) (Vehicle, error) {
	row := q.db.QueryRowContext(ctx, addVehicle,
		arg.CompanyID,
		arg.AssignedDriverID,
		arg.VehicleTypeID,
		arg.RegNo,
		arg.IsActive,
	)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.AssignedDriverID,
		&i.VehicleTypeID,
		&i.RegNo,
		&i.IsActive,
	)
	return i, err
}

const createVehicleTypes = `-- name: CreateVehicleTypes :one
INSERT INTO vehicle_types (id,name,max_vehicle_weight,max_vehicle_height,description) VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) ON CONFLICT(id) do update set name=EXCLUDED.name,max_vehicle_weight=EXCLUDED.max_vehicle_weight,max_vehicle_height=EXCLUDED.max_vehicle_height,description=EXCLUDED.description returning id, name, max_vehicle_weight, max_vehicle_height, description
`

type CreateVehicleTypesParams struct {
	ID               int32   `json:"id"`
	Name             string  `json:"name"`
	MaxVehicleWeight float64 `json:"max_vehicle_weight"`
	MaxVehicleHeight float64 `json:"max_vehicle_height"`
	Description      string  `json:"description"`
}

func (q *Queries) CreateVehicleTypes(ctx context.Context, arg CreateVehicleTypesParams) (VehicleType, error) {
	row := q.db.QueryRowContext(ctx, createVehicleTypes,
		arg.ID,
		arg.Name,
		arg.MaxVehicleWeight,
		arg.MaxVehicleHeight,
		arg.Description,
	)
	var i VehicleType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MaxVehicleWeight,
		&i.MaxVehicleHeight,
		&i.Description,
	)
	return i, err
}

const deleteVehicle = `-- name: DeleteVehicle :exec
delete from vehicles where id=$1
`

func (q *Queries) DeleteVehicle(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteVehicle, id)
	return err
}

const getAllVehicles = `-- name: GetAllVehicles :many
select vehicles.id, vehicles.company_id, vehicles.assigned_driver_id, vehicles.vehicle_type_id, vehicles.reg_no, vehicles.is_active,vehicle_types.description,vehicle_types.max_vehicle_height,vehicle_types.max_vehicle_weight,vehicle_types.name from vehicles inner join vehicle_types on vehicle_types.id=vehicles.vehicle_type_id where vehicles.company_id=$1
`

type GetAllVehiclesRow struct {
	ID               int32         `json:"id"`
	CompanyID        int32         `json:"company_id"`
	AssignedDriverID sql.NullInt32 `json:"assigned_driver_id"`
	VehicleTypeID    int32         `json:"vehicle_type_id"`
	RegNo            string        `json:"reg_no"`
	IsActive         bool          `json:"is_active"`
	Description      string        `json:"description"`
	MaxVehicleHeight float64       `json:"max_vehicle_height"`
	MaxVehicleWeight float64       `json:"max_vehicle_weight"`
	Name             string        `json:"name"`
}

func (q *Queries) GetAllVehicles(ctx context.Context, companyID int32) ([]GetAllVehiclesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllVehicles, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllVehiclesRow{}
	for rows.Next() {
		var i GetAllVehiclesRow
		if err := rows.Scan(
			&i.ID,
			&i.CompanyID,
			&i.AssignedDriverID,
			&i.VehicleTypeID,
			&i.RegNo,
			&i.IsActive,
			&i.Description,
			&i.MaxVehicleHeight,
			&i.MaxVehicleWeight,
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

const getDuplicateVehicle = `-- name: GetDuplicateVehicle :one
select count(*) from vehicles where lower(reg_no)=$1
`

func (q *Queries) GetDuplicateVehicle(ctx context.Context, regNo string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDuplicateVehicle, regNo)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDuplicateVehiclesWithoutID = `-- name: GetDuplicateVehiclesWithoutID :one
select count(*) from vehicles where lower(reg_no)=$1 and id!=$2
`

type GetDuplicateVehiclesWithoutIDParams struct {
	RegNo string `json:"reg_no"`
	ID    int32  `json:"id"`
}

func (q *Queries) GetDuplicateVehiclesWithoutID(ctx context.Context, arg GetDuplicateVehiclesWithoutIDParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDuplicateVehiclesWithoutID, arg.RegNo, arg.ID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getVehicleTypes = `-- name: GetVehicleTypes :many
select id, name, max_vehicle_weight, max_vehicle_height, description from vehicle_types
`

func (q *Queries) GetVehicleTypes(ctx context.Context) ([]VehicleType, error) {
	rows, err := q.db.QueryContext(ctx, getVehicleTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []VehicleType{}
	for rows.Next() {
		var i VehicleType
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MaxVehicleWeight,
			&i.MaxVehicleHeight,
			&i.Description,
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

const updateVehicle = `-- name: UpdateVehicle :exec
update vehicles set assigned_driver_id=$1,vehicle_type_id=$2,reg_no=$3,is_active=$4 where id=$5
`

type UpdateVehicleParams struct {
	AssignedDriverID sql.NullInt32 `json:"assigned_driver_id"`
	VehicleTypeID    int32         `json:"vehicle_type_id"`
	RegNo            string        `json:"reg_no"`
	IsActive         bool          `json:"is_active"`
	ID               int32         `json:"id"`
}

func (q *Queries) UpdateVehicle(ctx context.Context, arg UpdateVehicleParams) error {
	_, err := q.db.ExecContext(ctx, updateVehicle,
		arg.AssignedDriverID,
		arg.VehicleTypeID,
		arg.RegNo,
		arg.IsActive,
		arg.ID,
	)
	return err
}

const updateVehicleStatus = `-- name: UpdateVehicleStatus :exec
update vehicles set is_active=$1 where id =$2
`

type UpdateVehicleStatusParams struct {
	IsActive bool  `json:"is_active"`
	ID       int32 `json:"id"`
}

func (q *Queries) UpdateVehicleStatus(ctx context.Context, arg UpdateVehicleStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateVehicleStatus, arg.IsActive, arg.ID)
	return err
}
