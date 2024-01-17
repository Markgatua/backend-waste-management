// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: users.sql

package gen

import (
	"context"
	"database/sql"
	"time"
)

const createAdmin = `-- name: CreateAdmin :exec
insert into users (first_name,last_name,email,provider,role_id,email,password) VALUES($1,$2,$3,$4,$5,$6,$7)
`

type CreateAdminParams struct {
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Provider  sql.NullString `json:"provider"`
	RoleID    sql.NullInt32  `json:"role_id"`
	Email_2   sql.NullString `json:"email_2"`
	Password  sql.NullString `json:"password"`
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) error {
	_, err := q.db.ExecContext(ctx, createAdmin,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Provider,
		arg.RoleID,
		arg.Email_2,
		arg.Password,
	)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
select id, first_name, last_name, provider, role_id, user_company_id, email, password, avatar_url, user_type, is_active, calling_code, phone, phone_confirmed_at, confirmed_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, last_login, created_at, updated_at from users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Provider,
			&i.RoleID,
			&i.UserCompanyID,
			&i.Email,
			&i.Password,
			&i.AvatarUrl,
			&i.UserType,
			&i.IsActive,
			&i.CallingCode,
			&i.Phone,
			&i.PhoneConfirmedAt,
			&i.ConfirmedAt,
			&i.ConfirmationToken,
			&i.ConfirmationSentAt,
			&i.RecoveryToken,
			&i.RecoverySentAt,
			&i.LastLogin,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getUser = `-- name: GetUser :one
select id, first_name, last_name, provider, role_id, user_company_id, email, password, avatar_url, user_type, is_active, calling_code, phone, phone_confirmed_at, confirmed_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, last_login, created_at, updated_at from users where id=$1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Provider,
		&i.RoleID,
		&i.UserCompanyID,
		&i.Email,
		&i.Password,
		&i.AvatarUrl,
		&i.UserType,
		&i.IsActive,
		&i.CallingCode,
		&i.Phone,
		&i.PhoneConfirmedAt,
		&i.ConfirmedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.LastLogin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsersWithRole = `-- name: GetUsersWithRole :many
select users.id, users.first_name, users.last_name, users.provider, users.role_id, users.user_company_id, users.email, users.password, users.avatar_url, users.user_type, users.is_active, users.calling_code, users.phone, users.phone_confirmed_at, users.confirmed_at, users.confirmation_token, users.confirmation_sent_at, users.recovery_token, users.recovery_sent_at, users.last_login, users.created_at, users.updated_at, roles.name from users INNER JOIN roles ON users.role_id=roles.id
`

type GetUsersWithRoleRow struct {
	ID                 int32          `json:"id"`
	FirstName          sql.NullString `json:"first_name"`
	LastName           sql.NullString `json:"last_name"`
	Provider           sql.NullString `json:"provider"`
	RoleID             sql.NullInt32  `json:"role_id"`
	UserCompanyID      sql.NullInt32  `json:"user_company_id"`
	Email              sql.NullString `json:"email"`
	Password           sql.NullString `json:"password"`
	AvatarUrl          sql.NullString `json:"avatar_url"`
	UserType           sql.NullInt16  `json:"user_type"`
	IsActive           sql.NullBool   `json:"is_active"`
	CallingCode        sql.NullString `json:"calling_code"`
	Phone              sql.NullString `json:"phone"`
	PhoneConfirmedAt   sql.NullTime   `json:"phone_confirmed_at"`
	ConfirmedAt        sql.NullTime   `json:"confirmed_at"`
	ConfirmationToken  sql.NullString `json:"confirmation_token"`
	ConfirmationSentAt sql.NullTime   `json:"confirmation_sent_at"`
	RecoveryToken      sql.NullString `json:"recovery_token"`
	RecoverySentAt     sql.NullTime   `json:"recovery_sent_at"`
	LastLogin          sql.NullTime   `json:"last_login"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Name               string         `json:"name"`
}

func (q *Queries) GetUsersWithRole(ctx context.Context) ([]GetUsersWithRoleRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithRole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUsersWithRoleRow{}
	for rows.Next() {
		var i GetUsersWithRoleRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Provider,
			&i.RoleID,
			&i.UserCompanyID,
			&i.Email,
			&i.Password,
			&i.AvatarUrl,
			&i.UserType,
			&i.IsActive,
			&i.CallingCode,
			&i.Phone,
			&i.PhoneConfirmedAt,
			&i.ConfirmedAt,
			&i.ConfirmationToken,
			&i.ConfirmationSentAt,
			&i.RecoveryToken,
			&i.RecoverySentAt,
			&i.LastLogin,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateUser = `-- name: UpdateUser :exec
update users set first_name=$1,last_name=$2 where id=$3
`

type UpdateUserParams struct {
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser, arg.FirstName, arg.LastName, arg.ID)
	return err
}
