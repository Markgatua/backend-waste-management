// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package gen

import (
	"context"
)

type Querier interface {
	CreateCountry(ctx context.Context, arg CreateCountryParams) error
	CreatePermission(ctx context.Context, arg CreatePermissionParams) error
	// companies.sql
	GetAllCompanies(ctx context.Context) ([]Company, error)
	// regions.sql
	GetAllOrganizations(ctx context.Context) ([]Organization, error)
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	// waste_groups.sql
	GetAllWasteGroups(ctx context.Context) ([]WasteGroup, error)
	GetCompany(ctx context.Context, id int32) (Company, error)
	GetCountryBeCountryCode(ctx context.Context, countryCode string) ([]Country, error)
	GetDuplicateOrganization(ctx context.Context, arg GetDuplicateOrganizationParams) ([]Organization, error)
	GetOneWasteGroup(ctx context.Context, id int32) (WasteGroup, error)
	GetOrganization(ctx context.Context, id int32) (Organization, error)
	GetOrganizationCountWithNameAndCountry(ctx context.Context, arg GetOrganizationCountWithNameAndCountryParams) ([]Organization, error)
	GetUser(ctx context.Context, id int32) (User, error)
	GetUsersWithRole(ctx context.Context) ([]GetUsersWithRoleRow, error)
	InsertOrganization(ctx context.Context, arg InsertOrganizationParams) (Organization, error)
	InsertWasteGroup(ctx context.Context, arg InsertWasteGroupParams) (WasteGroup, error)
	UpdateCompanyStatus(ctx context.Context, arg UpdateCompanyStatusParams) error
	UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateWasteGroup(ctx context.Context, arg UpdateWasteGroupParams) error
}

var _ Querier = (*Queries)(nil)
