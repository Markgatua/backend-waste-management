// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package gen

import (
	"context"
	"database/sql"
)

type Querier interface {
	ActivateRole(ctx context.Context, roleID int32) error
	AddVehicle(ctx context.Context, arg AddVehicleParams) (Vehicle, error)
	AssignChampionToCollector(ctx context.Context, arg AssignChampionToCollectorParams) (ChampionAggregatorAssignment, error)
	AssignCollectorsToGreenChampion(ctx context.Context, arg AssignCollectorsToGreenChampionParams) (ChampionAggregatorAssignment, error)
	AssignPermission(ctx context.Context, arg AssignPermissionParams) error
	AssignPermissionToRole(ctx context.Context, arg AssignPermissionToRoleParams) error
	CancelCollectionRequest(ctx context.Context, id int32) error
	CheckSubCountiesDuplicate(ctx context.Context, name string) (int64, error)
	CollectionWeightTotals(ctx context.Context, producerID int32) (CollectionWeightTotalsRow, error)
	CompleteCollectionRequest(ctx context.Context, id int32) error
	ConfirmCollectionRequest(ctx context.Context, id int32) error
	CreateAggregatorWasteType(ctx context.Context, arg CreateAggregatorWasteTypeParams) (AggregatorWasteType, error)
	CreateBuyer(ctx context.Context, arg CreateBuyerParams) (Buyer, error)
	CreateCountry(ctx context.Context, arg CreateCountryParams) error
	CreateMainOrganizationAdmin(ctx context.Context, arg CreateMainOrganizationAdminParams) error
	CreatePermission(ctx context.Context, arg CreatePermissionParams) (Permission, error)
	CreatePurchase(ctx context.Context, arg CreatePurchaseParams) (Purchase, error)
	CreatePurchaseItem(ctx context.Context, arg CreatePurchaseItemParams) (PurchaseItem, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateSale(ctx context.Context, arg CreateSaleParams) (Sale, error)
	CreateSaleItem(ctx context.Context, arg CreateSaleItemParams) (SaleItem, error)
	CreateSupplier(ctx context.Context, arg CreateSupplierParams) (Supplier, error)
	CreateVehicleTypes(ctx context.Context, arg CreateVehicleTypesParams) (VehicleType, error)
	DeactivateRole(ctx context.Context, roleID int32) error
	DeleteAggregatorWasteTypes(ctx context.Context, aggregatorID int32) error
	DeleteBuyer(ctx context.Context, id int32) error
	DeleteChampionCollector(ctx context.Context, id int32) error
	DeleteCompany(ctx context.Context, id int32) error
	DeleteOrganization(ctx context.Context, id int32) error
	DeletePermissionByActions(ctx context.Context, actions []string) error
	DeletePermissionByIds(ctx context.Context, permissionIds []int32) error
	DeletePurchase(ctx context.Context, id int32) error
	DeleteRole(ctx context.Context, id int32) error
	DeleteSale(ctx context.Context, id int32) error
	DeleteSupplier(ctx context.Context, id int32) error
	DeleteVehicle(ctx context.Context, id int32) error
	DeleteWasteItemsForCollectionRequest(ctx context.Context, collectionRequestID int32) error
	DuplicateCounties(ctx context.Context, name string) (int64, error)
	DuplicatePickupTimeStamp(ctx context.Context, stamp string) (int64, error)
	GetAggregatorNewRequests(ctx context.Context, collectorID int32) ([]GetAggregatorNewRequestsRow, error)
	GetAggregatorWasteTypes(ctx context.Context, aggregatorID int32) ([]AggregatorWasteType, error)
	// companies.sql
	GetAllAggregators(ctx context.Context) ([]GetAllAggregatorsRow, error)
	GetAllCancelledCollectionRequests(ctx context.Context) ([]GetAllCancelledCollectionRequestsRow, error)
	// champion_aggregator_assignments.sql
	GetAllChampionCollectorsAssignments(ctx context.Context) ([]GetAllChampionCollectorsAssignmentsRow, error)
	GetAllChampionsForACollector(ctx context.Context, collectorID int32) ([]GetAllChampionsForACollectorRow, error)
	GetAllCollectionRequests(ctx context.Context) ([]GetAllCollectionRequestsRow, error)
	GetAllCollectionRequestsForACollector(ctx context.Context, collectorID int32) ([]GetAllCollectionRequestsForACollectorRow, error)
	GetAllCompletedCollectionRequests(ctx context.Context, status int32) ([]GetAllCompletedCollectionRequestsRow, error)
	GetAllCountries(ctx context.Context) ([]Country, error)
	GetAllGreenChampions(ctx context.Context) ([]GetAllGreenChampionsRow, error)
	GetAllMainOrganizationUsers(ctx context.Context) ([]GetAllMainOrganizationUsersRow, error)
	// regions.sql
	GetAllOrganizations(ctx context.Context) ([]GetAllOrganizationsRow, error)
	GetAllPendingCollectionRequests(ctx context.Context) ([]GetAllPendingCollectionRequestsRow, error)
	GetAllPendingConfirmationCollectionRequests(ctx context.Context) ([]GetAllPendingConfirmationCollectionRequestsRow, error)
	GetAllPermissionGroupedByModule(ctx context.Context) ([]Permission, error)
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetAllProducerCompletedCollectionRequests(ctx context.Context, producerID int32) ([]GetAllProducerCompletedCollectionRequestsRow, error)
	GetAllProducerPendingCollectionRequests(ctx context.Context, producerID int32) ([]GetAllProducerPendingCollectionRequestsRow, error)
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	GetAllVehicles(ctx context.Context, companyID int32) ([]GetAllVehiclesRow, error)
	// waste_types.sql
	GetAllWasteTypes(ctx context.Context) ([]GetAllWasteTypesRow, error)
	GetAssignedCollectorsToGreenChampion(ctx context.Context, championID int32) ([]ChampionAggregatorAssignment, error)
	GetBranchCount(ctx context.Context) ([]GetBranchCountRow, error)
	GetChildrenWasteTypes(ctx context.Context, parentID sql.NullInt32) ([]GetChildrenWasteTypesRow, error)
	GetCollectionStats(ctx context.Context, producerID int32) ([]GetCollectionStatsRow, error)
	GetCollectorsForGreenChampion(ctx context.Context, championID int32) ([]GetCollectorsForGreenChampionRow, error)
	GetCompany(ctx context.Context, id int32) (GetCompanyRow, error)
	GetCompanyUsers(ctx context.Context, userCompanyID sql.NullInt32) ([]GetCompanyUsersRow, error)
	GetCountryBeCountryCode(ctx context.Context, countryCode string) ([]Country, error)
	GetCountryByName(ctx context.Context, country string) (Country, error)
	GetDuplicateCompanies(ctx context.Context, arg GetDuplicateCompaniesParams) ([]Company, error)
	GetDuplicateCompaniesWithoutID(ctx context.Context, arg GetDuplicateCompaniesWithoutIDParams) ([]Company, error)
	GetDuplicateOrganization(ctx context.Context, arg GetDuplicateOrganizationParams) ([]Organization, error)
	GetDuplicateRoleHasPermission(ctx context.Context, arg GetDuplicateRoleHasPermissionParams) (int64, error)
	GetDuplicateVehicle(ctx context.Context, regNo string) (int64, error)
	GetDuplicateVehiclesWithoutID(ctx context.Context, arg GetDuplicateVehiclesWithoutIDParams) (int64, error)
	GetInventoryItem(ctx context.Context, arg GetInventoryItemParams) (Inventory, error)
	GetLatestCollection(ctx context.Context, id int32) (GetLatestCollectionRow, error)
	GetMainOrganization(ctx context.Context, organizationID string) ([]MainOrganization, error)
	GetMainOrganizationUser(ctx context.Context, id int32) (User, error)
	GetMainOrganizationUserByEmail(ctx context.Context, email sql.NullString) (User, error)
	GetMainWasteTypes(ctx context.Context) ([]GetMainWasteTypesRow, error)
	GetMyNotifications(ctx context.Context, userID int32) ([]Notification, error)
	GetOneWasteType(ctx context.Context, id int32) (GetOneWasteTypeRow, error)
	GetOrganization(ctx context.Context, id int32) (GetOrganizationRow, error)
	GetOrganizationCount(ctx context.Context) ([]GetOrganizationCountRow, error)
	GetOrganizationCountWithNameAndCountry(ctx context.Context, arg GetOrganizationCountWithNameAndCountryParams) ([]Organization, error)
	GetPermissionsForRoleID(ctx context.Context, roleID int32) ([]GetPermissionsForRoleIDRow, error)
	GetPickupTimeStamps(ctx context.Context) ([]PickupTimeStamp, error)
	GetProducerLatestCollectionId(ctx context.Context, producerID int32) (CollectionRequest, error)
	GetRole(ctx context.Context, id int32) (Role, error)
	// role_has_permissions.sql
	GetRolePermissions(ctx context.Context, roleID int32) ([]RoleHasPermission, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetSubCountiesForACounty(ctx context.Context, countyID int32) ([]SubCounty, error)
	GetTheCollectorForAChampion(ctx context.Context, championID int32) (GetTheCollectorForAChampionRow, error)
	GetUserWithEmailWithoutID(ctx context.Context, arg GetUserWithEmailWithoutIDParams) ([]User, error)
	GetUsersWasteType(ctx context.Context) ([]GetUsersWasteTypeRow, error)
	GetUsersWithRole(ctx context.Context) ([]GetUsersWithRoleRow, error)
	GetVehicleTypes(ctx context.Context) ([]VehicleType, error)
	GetWasteItemsProducerData(ctx context.Context, producerID int32) ([]GetWasteItemsProducerDataRow, error)
	// waste_items.sql
	InsertCollectionRequestWasteItem(ctx context.Context, arg InsertCollectionRequestWasteItemParams) (CollectionRequestWasteItem, error)
	InsertCompany(ctx context.Context, arg InsertCompanyParams) (Company, error)
	// counties.sql
	InsertCounties(ctx context.Context, name string) error
	InsertMainOrganization(ctx context.Context, arg InsertMainOrganizationParams) error
	// collection_requests.sql
	InsertNewCollectionRequest(ctx context.Context, arg InsertNewCollectionRequestParams) error
	// notifications.sql
	InsertNewNotificationRequest(ctx context.Context, arg InsertNewNotificationRequestParams) error
	InsertOrganization(ctx context.Context, arg InsertOrganizationParams) (Organization, error)
	// pickup_time_stamps.sql
	InsertPickupTimeStsmp(ctx context.Context, arg InsertPickupTimeStsmpParams) error
	// sub_counties.sql
	InsertSubcounties(ctx context.Context, arg InsertSubcountiesParams) error
	InsertToInventory(ctx context.Context, arg InsertToInventoryParams) error
	InsertWasteType(ctx context.Context, arg InsertWasteTypeParams) (WasteType, error)
	InventoryItemCount(ctx context.Context, arg InventoryItemCountParams) (int64, error)
	MakeCashPayment(ctx context.Context, arg MakeCashPaymentParams) (SaleTransaction, error)
	MakePurchaseCashPayment(ctx context.Context, arg MakePurchaseCashPaymentParams) (PurchaseTransaction, error)
	RemoveAggrigatorsAssignedFromGreenChampions(ctx context.Context, championID int32) error
	RemovePermissionsFromRole(ctx context.Context, arg RemovePermissionsFromRoleParams) error
	RemoveRolePermissions(ctx context.Context, roleID int32) error
	RevokePermission(ctx context.Context, arg RevokePermissionParams) error
	RoleExists(ctx context.Context, name string) (int64, error)
	SetBuyerActiveInactiveStatus(ctx context.Context, arg SetBuyerActiveInactiveStatusParams) error
	SetPickupTimesForGreenChampion(ctx context.Context, arg SetPickupTimesForGreenChampionParams) error
	SetSupplierActiveInactiveStatus(ctx context.Context, arg SetSupplierActiveInactiveStatusParams) error
	UpdateBuyer(ctx context.Context, arg UpdateBuyerParams) error
	UpdateChampionCollector(ctx context.Context, arg UpdateChampionCollectorParams) error
	UpdateCollectionRequest(ctx context.Context, arg UpdateCollectionRequestParams) error
	UpdateCompany(ctx context.Context, arg UpdateCompanyParams) error
	UpdateCompanyStatus(ctx context.Context, arg UpdateCompanyStatusParams) error
	UpdateInventoryItem(ctx context.Context, arg UpdateInventoryItemParams) error
	// ttnm_organization.sql
	UpdateMainOrganizationProfile(ctx context.Context, arg UpdateMainOrganizationProfileParams) error
	UpdateMainOrganizationUser(ctx context.Context, arg UpdateMainOrganizationUserParams) error
	UpdateNotificationStatus(ctx context.Context, arg UpdateNotificationStatusParams) error
	UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) error
	UpdateOrganizationIsActive(ctx context.Context, arg UpdateOrganizationIsActiveParams) error
	UpdatePickupTimeStamp(ctx context.Context, arg UpdatePickupTimeStampParams) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) error
	UpdateSupplier(ctx context.Context, arg UpdateSupplierParams) error
	UpdateUserFirstNameLastNameEmailRoleAndUserType(ctx context.Context, arg UpdateUserFirstNameLastNameEmailRoleAndUserTypeParams) error
	UpdateUserFirstNameLastNameEmailRoleUserTypeAndPassword(ctx context.Context, arg UpdateUserFirstNameLastNameEmailRoleUserTypeAndPasswordParams) error
	UpdateUserIsActive(ctx context.Context, arg UpdateUserIsActiveParams) error
	UpdateVehicle(ctx context.Context, arg UpdateVehicleParams) error
	UpdateVehicleStatus(ctx context.Context, arg UpdateVehicleStatusParams) error
	UpdateWasteType(ctx context.Context, arg UpdateWasteTypeParams) error
	ViewCounties(ctx context.Context) ([]County, error)
	ViewSubCounties(ctx context.Context) ([]SubCounty, error)
}

var _ Querier = (*Queries)(nil)
