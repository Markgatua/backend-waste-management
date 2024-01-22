// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package gen

import (
	"context"
	"database/sql"
)

type Querier interface {
	ActivateRole(ctx context.Context, roleID int32) error
	AssignChampionToCollector(ctx context.Context, arg AssignChampionToCollectorParams) (ChampionAggregatorAssignment, error)
	AssignPermission(ctx context.Context, arg AssignPermissionParams) error
	AssignPermissionToRole(ctx context.Context, arg AssignPermissionToRoleParams) error
	CancelCollectionRequest(ctx context.Context, arg CancelCollectionRequestParams) error
	CheckSubCountiesDuplicate(ctx context.Context, name string) (int64, error)
	CollectionWeightTotals(ctx context.Context, producerID int32) (CollectionWeightTotalsRow, error)
	ConfirmCollectionRequest(ctx context.Context, arg ConfirmCollectionRequestParams) error
	CreateCountry(ctx context.Context, arg CreateCountryParams) error
	CreateMainOrganizationAdmin(ctx context.Context, arg CreateMainOrganizationAdminParams) error
	CreatePermission(ctx context.Context, arg CreatePermissionParams) (Permission, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	DeactivateRole(ctx context.Context, roleID int32) error
	DeleteChampionCollector(ctx context.Context, id int32) error
	DeleteCompany(ctx context.Context, id int32) error
	DeleteOrganization(ctx context.Context, id int32) error
	DeletePermissionByActions(ctx context.Context, actions []string) error
	DeletePermissionByIds(ctx context.Context, permissionIds []int32) error
	DeleteRole(ctx context.Context, id int32) error
	DuplicateCounties(ctx context.Context, name string) (int64, error)
	GetAllCancelledCollectionRequests(ctx context.Context, cancelled sql.NullBool) ([]GetAllCancelledCollectionRequestsRow, error)
	// champion_aggregator_assignments.sql
	GetAllChampionCollectorsAssignments(ctx context.Context) ([]GetAllChampionCollectorsAssignmentsRow, error)
	GetAllChampionsForACollector(ctx context.Context, collectorID sql.NullInt32) ([]GetAllChampionsForACollectorRow, error)
	GetAllCollectionRequests(ctx context.Context) ([]GetAllCollectionRequestsRow, error)
	GetAllCollectionRequestsForACollector(ctx context.Context, collectorID int32) ([]GetAllCollectionRequestsForACollectorRow, error)
	// companies.sql
	GetAllCompanies(ctx context.Context) ([]GetAllCompaniesRow, error)
	GetAllCompletedCollectionRequests(ctx context.Context, status sql.NullBool) ([]GetAllCompletedCollectionRequestsRow, error)
	GetAllCountries(ctx context.Context) ([]Country, error)
	GetAllMainOrganizationUsers(ctx context.Context) ([]GetAllMainOrganizationUsersRow, error)
	// regions.sql
	GetAllOrganizations(ctx context.Context) ([]GetAllOrganizationsRow, error)
	GetAllPendingCollectionRequests(ctx context.Context, confirmed sql.NullBool) ([]GetAllPendingCollectionRequestsRow, error)
	GetAllPendingConfirmationCollectionRequests(ctx context.Context, confirmed sql.NullBool) ([]GetAllPendingConfirmationCollectionRequestsRow, error)
	GetAllPermissionGroupedByModule(ctx context.Context) ([]Permission, error)
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	// waste_types.sql
	GetAllWasteTypes(ctx context.Context) ([]GetAllWasteTypesRow, error)
	GetAssignedCollectorsToGreenChampion(ctx context.Context, championID sql.NullInt32) ([]ChampionAggregatorAssignment, error)
	GetCompany(ctx context.Context, id int32) (GetCompanyRow, error)
	GetCompanyUsers(ctx context.Context, userCompanyID sql.NullInt32) ([]GetCompanyUsersRow, error)
	GetCountryBeCountryCode(ctx context.Context, countryCode string) ([]Country, error)
	GetDuplicateCompanies(ctx context.Context, arg GetDuplicateCompaniesParams) ([]Company, error)
	GetDuplicateCompaniesWithoutID(ctx context.Context, arg GetDuplicateCompaniesWithoutIDParams) ([]Company, error)
	GetDuplicateOrganization(ctx context.Context, arg GetDuplicateOrganizationParams) ([]Organization, error)
	GetDuplicateRoleHasPermission(ctx context.Context, arg GetDuplicateRoleHasPermissionParams) (int64, error)
	GetMainOrganization(ctx context.Context, organizationID string) ([]MainOrganization, error)
	GetMainOrganizationUser(ctx context.Context, id int32) (User, error)
	GetMainOrganizationUserByEmail(ctx context.Context, email sql.NullString) (User, error)
	GetMyNotifications(ctx context.Context, userID int32) ([]Notification, error)
	GetOneWasteType(ctx context.Context, id int32) (WasteType, error)
	GetOrganization(ctx context.Context, id int32) (GetOrganizationRow, error)
	GetOrganizationCountWithNameAndCountry(ctx context.Context, arg GetOrganizationCountWithNameAndCountryParams) ([]Organization, error)
	GetPermissionsForRoleID(ctx context.Context, roleID int32) ([]GetPermissionsForRoleIDRow, error)
	GetRole(ctx context.Context, id int32) (Role, error)
	// role_has_permissions.sql
	GetRolePermissions(ctx context.Context, roleID int32) ([]RoleHasPermission, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetSubCountiesForACounty(ctx context.Context, countyID int32) ([]SubCounty, error)
	GetTheCollectorForAChampion(ctx context.Context, championID sql.NullInt32) (GetTheCollectorForAChampionRow, error)
	GetUsersWasteType(ctx context.Context) ([]WasteType, error)
	GetUsersWithRole(ctx context.Context) ([]GetUsersWithRoleRow, error)
	InsertCompany(ctx context.Context, arg InsertCompanyParams) (Company, error)
	// counties.sql
	InsertCounties(ctx context.Context, name string) error
	InsertMainOrganization(ctx context.Context, arg InsertMainOrganizationParams) error
	// collection_requests.sql
	InsertNewCollectionRequest(ctx context.Context, arg InsertNewCollectionRequestParams) error
	// notifications.sql
	InsertNewNotificationRequest(ctx context.Context, arg InsertNewNotificationRequestParams) error
	InsertOrganization(ctx context.Context, arg InsertOrganizationParams) (Organization, error)
	// sub_counties.sql
	InsertSubcounties(ctx context.Context, arg InsertSubcountiesParams) error
	// waste_items.sql
	InsertWasteItem(ctx context.Context, arg InsertWasteItemParams) (WasteItem, error)
	InsertWasteType(ctx context.Context, arg InsertWasteTypeParams) (WasteType, error)
	RemovePermissionsFromRole(ctx context.Context, arg RemovePermissionsFromRoleParams) error
	RemoveRolePermissions(ctx context.Context, roleID int32) error
	RevokePermission(ctx context.Context, arg RevokePermissionParams) error
	RoleExists(ctx context.Context, name string) (int64, error)
	UpdateChampionCollector(ctx context.Context, arg UpdateChampionCollectorParams) error
	UpdateCollectionRequest(ctx context.Context, arg UpdateCollectionRequestParams) error
	UpdateCompany(ctx context.Context, arg UpdateCompanyParams) error
	UpdateCompanyStatus(ctx context.Context, arg UpdateCompanyStatusParams) error
	// ttnm_organization.sql
	UpdateMainOrganizationProfile(ctx context.Context, arg UpdateMainOrganizationProfileParams) error
	UpdateMainOrganizationUser(ctx context.Context, arg UpdateMainOrganizationUserParams) error
	UpdateNotificationStatus(ctx context.Context, arg UpdateNotificationStatusParams) error
	UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) error
	UpdateUserIsActive(ctx context.Context, arg UpdateUserIsActiveParams) error
	UpdateWasteType(ctx context.Context, arg UpdateWasteTypeParams) error
	ViewCounties(ctx context.Context) ([]County, error)
	ViewSubCounties(ctx context.Context) ([]SubCounty, error)
}

var _ Querier = (*Queries)(nil)
