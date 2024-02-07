// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package gen

import (
	"context"
	"database/sql"
	"time"
)

const createMainOrganizationAdmin = `-- name: CreateMainOrganizationAdmin :exec
insert into
    users (
        first_name, last_name, email, provider, role_id, password, confirmed_at, is_main_organization_user
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8
    )
`

type CreateMainOrganizationAdminParams struct {
	FirstName              sql.NullString `json:"first_name"`
	LastName               sql.NullString `json:"last_name"`
	Email                  sql.NullString `json:"email"`
	Provider               sql.NullString `json:"provider"`
	RoleID                 sql.NullInt32  `json:"role_id"`
	Password               sql.NullString `json:"password"`
	ConfirmedAt            sql.NullTime   `json:"confirmed_at"`
	IsMainOrganizationUser bool           `json:"is_main_organization_user"`
}

func (q *Queries) CreateMainOrganizationAdmin(ctx context.Context, arg CreateMainOrganizationAdminParams) error {
	_, err := q.db.ExecContext(ctx, createMainOrganizationAdmin,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Provider,
		arg.RoleID,
		arg.Password,
		arg.ConfirmedAt,
		arg.IsMainOrganizationUser,
	)
	return err
}

const getAllMainOrganizationUsers = `-- name: GetAllMainOrganizationUsers :many
select
    users.id,
    users.first_name,
    users.last_name,
    users.email,
    users.avatar_url,
    users.calling_code,
    users.phone,
    users.is_active,
    roles.name as role_name,
    roles.id as role_id
from users
    inner join roles on users.role_id = roles.id
where
    users.email not ilike 'superadmin@admin.com'
    and users.is_main_organization_user = true
`

type GetAllMainOrganizationUsersRow struct {
	ID          int32          `json:"id"`
	FirstName   sql.NullString `json:"first_name"`
	LastName    sql.NullString `json:"last_name"`
	Email       sql.NullString `json:"email"`
	AvatarUrl   sql.NullString `json:"avatar_url"`
	CallingCode sql.NullString `json:"calling_code"`
	Phone       sql.NullString `json:"phone"`
	IsActive    sql.NullBool   `json:"is_active"`
	RoleName    string         `json:"role_name"`
	RoleID      int32          `json:"role_id"`
}

func (q *Queries) GetAllMainOrganizationUsers(ctx context.Context) ([]GetAllMainOrganizationUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMainOrganizationUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllMainOrganizationUsersRow{}
	for rows.Next() {
		var i GetAllMainOrganizationUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.AvatarUrl,
			&i.CallingCode,
			&i.Phone,
			&i.IsActive,
			&i.RoleName,
			&i.RoleID,
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

const getAllUsers = `-- name: GetAllUsers :many
select
    users.id,
    users.first_name,
    users.last_name,
    users.email,
    users.avatar_url,
    users.calling_code,
    users.phone,
    users.is_active,
    roles.name as role_name,
    roles.id as role_id
from users
    inner join roles on users.role_id = roles.id
where
    users.email not ilike 'superadmin@admin.com'
    and users.is_main_organization_user = false
`

type GetAllUsersRow struct {
	ID          int32          `json:"id"`
	FirstName   sql.NullString `json:"first_name"`
	LastName    sql.NullString `json:"last_name"`
	Email       sql.NullString `json:"email"`
	AvatarUrl   sql.NullString `json:"avatar_url"`
	CallingCode sql.NullString `json:"calling_code"`
	Phone       sql.NullString `json:"phone"`
	IsActive    sql.NullBool   `json:"is_active"`
	RoleName    string         `json:"role_name"`
	RoleID      int32          `json:"role_id"`
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllUsersRow{}
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.AvatarUrl,
			&i.CallingCode,
			&i.Phone,
			&i.IsActive,
			&i.RoleName,
			&i.RoleID,
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

const getCompanyUsers = `-- name: GetCompanyUsers :many
select users.id, users.first_name, users.last_name, users.provider, users.role_id, users.user_company_id, users.user_organization_id, users.is_main_organization_user, users.is_organization_super_admin, users.is_company_super_admin, users.email, users.password, users.avatar_url, users.user_type, users.is_active, users.calling_code, users.phone, users.phone_confirmed_at, users.confirmed_at, users.confirmation_token, users.confirmation_sent_at, users.recovery_token, users.recovery_sent_at, users.last_login, users.created_at, users.updated_at,companies.name as company,companies.location as companyLocation FROM users
LEFT JOIN
  companies ON companies.id = users.user_company_id
 where user_company_id = $1
`

type GetCompanyUsersRow struct {
	ID                       int32          `json:"id"`
	FirstName                sql.NullString `json:"first_name"`
	LastName                 sql.NullString `json:"last_name"`
	Provider                 sql.NullString `json:"provider"`
	RoleID                   sql.NullInt32  `json:"role_id"`
	UserCompanyID            sql.NullInt32  `json:"user_company_id"`
	UserOrganizationID       sql.NullInt32  `json:"user_organization_id"`
	IsMainOrganizationUser   bool           `json:"is_main_organization_user"`
	IsOrganizationSuperAdmin bool           `json:"is_organization_super_admin"`
	IsCompanySuperAdmin      bool           `json:"is_company_super_admin"`
	Email                    sql.NullString `json:"email"`
	Password                 sql.NullString `json:"password"`
	AvatarUrl                sql.NullString `json:"avatar_url"`
	UserType                 sql.NullInt16  `json:"user_type"`
	IsActive                 sql.NullBool   `json:"is_active"`
	CallingCode              sql.NullString `json:"calling_code"`
	Phone                    sql.NullString `json:"phone"`
	PhoneConfirmedAt         sql.NullTime   `json:"phone_confirmed_at"`
	ConfirmedAt              sql.NullTime   `json:"confirmed_at"`
	ConfirmationToken        sql.NullString `json:"confirmation_token"`
	ConfirmationSentAt       sql.NullTime   `json:"confirmation_sent_at"`
	RecoveryToken            sql.NullString `json:"recovery_token"`
	RecoverySentAt           sql.NullTime   `json:"recovery_sent_at"`
	LastLogin                sql.NullTime   `json:"last_login"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	Company                  sql.NullString `json:"company"`
	Companylocation          sql.NullString `json:"companylocation"`
}

func (q *Queries) GetCompanyUsers(ctx context.Context, userCompanyID sql.NullInt32) ([]GetCompanyUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getCompanyUsers, userCompanyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCompanyUsersRow{}
	for rows.Next() {
		var i GetCompanyUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Provider,
			&i.RoleID,
			&i.UserCompanyID,
			&i.UserOrganizationID,
			&i.IsMainOrganizationUser,
			&i.IsOrganizationSuperAdmin,
			&i.IsCompanySuperAdmin,
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
			&i.Company,
			&i.Companylocation,
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

const getMainOrganizationUser = `-- name: GetMainOrganizationUser :one
select id, first_name, last_name, provider, role_id, user_company_id, user_organization_id, is_main_organization_user, is_organization_super_admin, is_company_super_admin, email, password, avatar_url, user_type, is_active, calling_code, phone, phone_confirmed_at, confirmed_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, last_login, created_at, updated_at from users where id = $1
`

func (q *Queries) GetMainOrganizationUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getMainOrganizationUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Provider,
		&i.RoleID,
		&i.UserCompanyID,
		&i.UserOrganizationID,
		&i.IsMainOrganizationUser,
		&i.IsOrganizationSuperAdmin,
		&i.IsCompanySuperAdmin,
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

const getMainOrganizationUserByEmail = `-- name: GetMainOrganizationUserByEmail :one
select id, first_name, last_name, provider, role_id, user_company_id, user_organization_id, is_main_organization_user, is_organization_super_admin, is_company_super_admin, email, password, avatar_url, user_type, is_active, calling_code, phone, phone_confirmed_at, confirmed_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, last_login, created_at, updated_at from users where email = $1
`

func (q *Queries) GetMainOrganizationUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getMainOrganizationUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Provider,
		&i.RoleID,
		&i.UserCompanyID,
		&i.UserOrganizationID,
		&i.IsMainOrganizationUser,
		&i.IsOrganizationSuperAdmin,
		&i.IsCompanySuperAdmin,
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

const getUserWithEmailWithoutID = `-- name: GetUserWithEmailWithoutID :many
select id, first_name, last_name, provider, role_id, user_company_id, user_organization_id, is_main_organization_user, is_organization_super_admin, is_company_super_admin, email, password, avatar_url, user_type, is_active, calling_code, phone, phone_confirmed_at, confirmed_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, last_login, created_at, updated_at from users where email = $1 and id != $2
`

type GetUserWithEmailWithoutIDParams struct {
	Email sql.NullString `json:"email"`
	ID    int32          `json:"id"`
}

func (q *Queries) GetUserWithEmailWithoutID(ctx context.Context, arg GetUserWithEmailWithoutIDParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUserWithEmailWithoutID, arg.Email, arg.ID)
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
			&i.UserOrganizationID,
			&i.IsMainOrganizationUser,
			&i.IsOrganizationSuperAdmin,
			&i.IsCompanySuperAdmin,
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

const getUsersWithRole = `-- name: GetUsersWithRole :many
select users.id, users.first_name, users.last_name, users.provider, users.role_id, users.user_company_id, users.user_organization_id, users.is_main_organization_user, users.is_organization_super_admin, users.is_company_super_admin, users.email, users.password, users.avatar_url, users.user_type, users.is_active, users.calling_code, users.phone, users.phone_confirmed_at, users.confirmed_at, users.confirmation_token, users.confirmation_sent_at, users.recovery_token, users.recovery_sent_at, users.last_login, users.created_at, users.updated_at, roles.name
from users
    INNER JOIN roles ON users.role_id = roles.id
`

type GetUsersWithRoleRow struct {
	ID                       int32          `json:"id"`
	FirstName                sql.NullString `json:"first_name"`
	LastName                 sql.NullString `json:"last_name"`
	Provider                 sql.NullString `json:"provider"`
	RoleID                   sql.NullInt32  `json:"role_id"`
	UserCompanyID            sql.NullInt32  `json:"user_company_id"`
	UserOrganizationID       sql.NullInt32  `json:"user_organization_id"`
	IsMainOrganizationUser   bool           `json:"is_main_organization_user"`
	IsOrganizationSuperAdmin bool           `json:"is_organization_super_admin"`
	IsCompanySuperAdmin      bool           `json:"is_company_super_admin"`
	Email                    sql.NullString `json:"email"`
	Password                 sql.NullString `json:"password"`
	AvatarUrl                sql.NullString `json:"avatar_url"`
	UserType                 sql.NullInt16  `json:"user_type"`
	IsActive                 sql.NullBool   `json:"is_active"`
	CallingCode              sql.NullString `json:"calling_code"`
	Phone                    sql.NullString `json:"phone"`
	PhoneConfirmedAt         sql.NullTime   `json:"phone_confirmed_at"`
	ConfirmedAt              sql.NullTime   `json:"confirmed_at"`
	ConfirmationToken        sql.NullString `json:"confirmation_token"`
	ConfirmationSentAt       sql.NullTime   `json:"confirmation_sent_at"`
	RecoveryToken            sql.NullString `json:"recovery_token"`
	RecoverySentAt           sql.NullTime   `json:"recovery_sent_at"`
	LastLogin                sql.NullTime   `json:"last_login"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	Name                     string         `json:"name"`
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
			&i.UserOrganizationID,
			&i.IsMainOrganizationUser,
			&i.IsOrganizationSuperAdmin,
			&i.IsCompanySuperAdmin,
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

const updateMainOrganizationUser = `-- name: UpdateMainOrganizationUser :exec
update users set first_name = $1, last_name = $2 where id = $3
`

type UpdateMainOrganizationUserParams struct {
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateMainOrganizationUser(ctx context.Context, arg UpdateMainOrganizationUserParams) error {
	_, err := q.db.ExecContext(ctx, updateMainOrganizationUser, arg.FirstName, arg.LastName, arg.ID)
	return err
}

const updateUserFirstNameLastNameEmailRoleAndUserType = `-- name: UpdateUserFirstNameLastNameEmailRoleAndUserType :exec
update users
set
    email = $1,
    role_id = $2,
    user_type = $3,
    first_name=$4,
    last_name=$5
where
    id = $6
`

type UpdateUserFirstNameLastNameEmailRoleAndUserTypeParams struct {
	Email     sql.NullString `json:"email"`
	RoleID    sql.NullInt32  `json:"role_id"`
	UserType  sql.NullInt16  `json:"user_type"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateUserFirstNameLastNameEmailRoleAndUserType(ctx context.Context, arg UpdateUserFirstNameLastNameEmailRoleAndUserTypeParams) error {
	_, err := q.db.ExecContext(ctx, updateUserFirstNameLastNameEmailRoleAndUserType,
		arg.Email,
		arg.RoleID,
		arg.UserType,
		arg.FirstName,
		arg.LastName,
		arg.ID,
	)
	return err
}

const updateUserFirstNameLastNameEmailRoleUserTypeAndPassword = `-- name: UpdateUserFirstNameLastNameEmailRoleUserTypeAndPassword :exec
update users
set
    email = $1,
    role_id = $2,
    password = $3,
    user_type = $4,
    first_name=$5,
    last_name=$6
where
    id = $7
`

type UpdateUserFirstNameLastNameEmailRoleUserTypeAndPasswordParams struct {
	Email     sql.NullString `json:"email"`
	RoleID    sql.NullInt32  `json:"role_id"`
	Password  sql.NullString `json:"password"`
	UserType  sql.NullInt16  `json:"user_type"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateUserFirstNameLastNameEmailRoleUserTypeAndPassword(ctx context.Context, arg UpdateUserFirstNameLastNameEmailRoleUserTypeAndPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserFirstNameLastNameEmailRoleUserTypeAndPassword,
		arg.Email,
		arg.RoleID,
		arg.Password,
		arg.UserType,
		arg.FirstName,
		arg.LastName,
		arg.ID,
	)
	return err
}

const updateUserIsActive = `-- name: UpdateUserIsActive :exec
update users set is_active = $1 where id = $2
`

type UpdateUserIsActiveParams struct {
	IsActive sql.NullBool `json:"is_active"`
	ID       int32        `json:"id"`
}

func (q *Queries) UpdateUserIsActive(ctx context.Context, arg UpdateUserIsActiveParams) error {
	_, err := q.db.ExecContext(ctx, updateUserIsActive, arg.IsActive, arg.ID)
	return err
}
