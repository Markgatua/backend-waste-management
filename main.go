package main

import (
	"fmt"
	_ "fmt"
	"os"
	_ "time"
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/controllers"
	"ttnmwastemanagementsystem/database/seeder"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	var instructions string = `
	   Instructions
	   ------------
	   1) To run seeder run "go run main.go seeder"
	   2) To run master table migration run "go run main.go master_table_migration"
	   3) To run the system run "go run main.go program"
	`

	if len(os.Args) == 1 {
		fmt.Print(instructions)
	} else {
		args := os.Args[1]
		cache2go.Cache("scms_cache")

		if args == "seeder" {
			seeder.Run()
		} else if args == "master_table_migration" {
			//migrations.RunMasterTableMigration()
		} else if args == "program" {
			runProgram()
		} else {
			fmt.Println(instructions)
			//panic(fmt.Sprint("Unkown arg [", args, "]"))
		}
	}
}
func runProgram() {
	router := gin.Default()
	configs.InitEnvConfigs(".")

	gen.LoadRepo()

	usersController := controllers.UsersController{}
	wasteTypesController := controllers.WasteTypesController{}
	organzationController := controllers.OrgnizationController{}
	geoController := controllers.GeoController{}
	championCollectorController := controllers.ChampionCollectorController{}
	// rolesController := controllers.RolesController{}
	// rolespermissions := controllers.RoleAndPermissionsController{}
	ttnmOrganizationController := controllers.TtnmOrganizationController{}
	requestCollectionController := controllers.RequestCollectionController{}
	wasteItemsController := controllers.WasteItemsController{}

	router.LoadHTMLGlob("templates/**/*")

	//router.Use(cors.Default())
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"PUT", "POST", "DELETE", "OPTIONS", "GET", "PATCH"},
		AllowCredentials: true,

		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	//router.Use(cors.Default())
	router.Use(c)
	authController := controllers.AuthController{}
	auth := router.Group("/auth")
	{
		auth.GET("/challenge/password_reset_success/email/web", authController.PassWordResetSuccess)
		auth.POST("/challenge/submit_new_password/email/web", authController.SubmitNewPassword)
		auth.GET("/challenge/enter_new_password/web", authController.EnterNewPassword)
		auth.POST("/reset_password/email/web", authController.ResetPassword)

		auth.POST("/register/email/:organization", authController.RegisterUserEmail)
		auth.GET("/challenge/verify_email/web", authController.VerifyEmail)
		auth.POST("/challenge/send_email_verfication/web", authController.SendVerificationMail)
		auth.POST("/login/email", authController.LoginEmail)
		auth.POST("/register/phone/update_profile", authController.RegisterPhoneUpdateUserDetails)

		auth.POST("/challenge/register/send_otp_code_phone", authController.RegisterPhoneSendOTPCode)
		auth.POST("/login/phone", authController.LoginPhone)
		auth.POST("/challenge/register/verify_otp_code_phone", authController.RegisterVerifyOTPCodePhone)
		auth.POST("/challenge/register/verify_otp_code_phone_create_wallet", authController.RegisterVerifyOTPCodePhoneAndCreateWallet)

		auth.POST("/challenge/forgot_pin/send_otp_phone", authController.ForgotPinSendOTPPhone)
		auth.POST("/challenge/forgot_pin/verify_otp_phone", authController.ForgotPinVerifyOTPPhone)
		auth.POST("/challenge/forgot_pin/enter_new_pin", authController.ForgotPinEnterNewPin)

		auth.POST("/reset_password/email/mobile", authController.ResetPasswordApi)
		auth.POST("/reset_password/email/mobile/reset", authController.SubmitNewPasswordApi)
		auth.POST("/reset_password/phone/mobile/reset", authController.PasswordResetAndVerifyOTPPhone)

	}

	// apiGroup := router.Group("/api:BVaDN9hl")
	router.GET("/uploads/:file", controllers.FileController{}.GetFile)
	router.GET("/flag/:file", controllers.FileController{}.GetFlag)

	//presets ----------------------------------------------------------|
	router.GET("/presets", controllers.PresetController{}.GetPresetValue)
	//------------------------------------------------------------------|

	router.Use(middlewares.JwtAuthMiddleware())
	// router.Use(middlewares.PermissionMiddleware())

	//---------------------------    Files ------------------------------------------------------
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload_files", middlewares.PermissionBlockerMiddleware("upload_files"), controllers.FileController{}.UploadFiles)
	//-------------------------------------------------------------------------------------------

	//stats --------
	router.GET("/main_organization_stats", middlewares.PermissionBlockerMiddleware("view_stats"), controllers.StatsController{}.GetMainOrganizationStats)
	//stats -------

	router.GET("/users", middlewares.PermissionBlockerMiddleware("view_user"), usersController.GetAllUsers)
	//main organizations is
	router.GET("/users/main_organization", middlewares.PermissionBlockerMiddleware("view_user"), usersController.GetAllMainOrganizationUsers)
	router.PUT("/users/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_user"), usersController.SetActiveInActiveStatus)
	router.PUT("/users/edit/email/:organization", middlewares.PermissionBlockerMiddleware("edit_user"), authController.EditUserEmail)
	router.PUT("/user/update_password", middlewares.PermissionBlockerMiddleware("update_user_password"), authController.UpdateUserPassword)

	// router.POST("update/user",usersController.UpdateUSer)
	// router.GET("/users/roles",usersController.GetUsersWithRole)
	router.GET("/user/:id/main_organization", middlewares.PermissionBlockerMiddleware("view_user"), usersController.GetMainOrganizationUser)
	router.GET("/user/:id", middlewares.PermissionBlockerMiddleware("view_user"), usersController.GetUser)

	router.GET("/company/users/:id", usersController.GetCompanyUsers)

	//---------------------------countries-------------------------------------------------------
	router.GET("countries", geoController.GetAllCountries)
	//-------------------------------------------------------------------------------------------

	//---------------------------buyer----------------------------------------------------
	router.POST("buyer/add", middlewares.PermissionBlockerMiddleware("add_buyer"), controllers.AggregatorController{}.AddBuyer)
	router.PUT("buyer/update", middlewares.PermissionBlockerMiddleware("edit_buyer"), controllers.AggregatorController{}.UpdateBuyer)
	router.GET("buyers", middlewares.PermissionBlockerMiddleware("view_buyer"), controllers.AggregatorController{}.GetBuyers)
	router.PUT("/buyer/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_buyer"), controllers.AggregatorController{}.SetBuyerActiveInActiveStatus)

	router.DELETE("buyer/delete/:id", middlewares.PermissionBlockerMiddleware("delete_buyer"), controllers.AggregatorController{}.DeleteBuyer)
	//-------------------------------------------------------------------------------------------

	//---------------------------suppliers----------------------------------------------------
	router.POST("supplier/add", middlewares.PermissionBlockerMiddleware("add_supplier"), controllers.AggregatorController{}.AddSupplier)
	router.PUT("supplier/update", middlewares.PermissionBlockerMiddleware("edit_supplier"), controllers.AggregatorController{}.UpdateSupplier)
	router.GET("suppliers", middlewares.PermissionBlockerMiddleware("view_supplier"), controllers.AggregatorController{}.GetSuppliers)
	router.DELETE("supplier/delete/:id", middlewares.PermissionBlockerMiddleware("delete_supplier"), controllers.AggregatorController{}.DeleteSupplier)
	router.PUT("/supplier/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_supplier"), controllers.AggregatorController{}.SetSupplierActiveInActiveStatus)

	//-------------------------------------------------------------------------------------------

	//---------------------------organization----------------------------------------------------
	router.POST("organization/add", middlewares.PermissionBlockerMiddleware("add_organization"), organzationController.InsertOrganization)
	router.PUT("organization/update", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.UpdateOrganization)
	router.GET("organizations", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetAllOrganizations)
	router.PUT("organization/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.SetActiveInActiveStatus)

	router.DELETE("organization/delete/:id", middlewares.PermissionBlockerMiddleware("delete_organization"), organzationController.DeleteOrganization)
	router.GET("organization/:id", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetOrganization)
	//-------------------------------------------------------------------------------------------

	//---------------------------Sell-------------------------------------------------------------
	router.POST("aggregator/sell_waste_to_buyer", middlewares.PermissionBlockerMiddleware("sell"), controllers.AggregatorController{}.SellWasteToBuyer)
	router.GET("aggregator/sales", middlewares.PermissionBlockerMiddleware("view_sale_history"), controllers.AggregatorController{}.GetSales)
	//--------------------------------------------------------------------------------------------

	//---------------------------Purchases-------------------------------------------------------------
	router.POST("aggregator/purchase_waste_from_supplier", middlewares.PermissionBlockerMiddleware("purchase"), controllers.AggregatorController{}.PurchaseWasteFromSupplier)
	router.GET("aggregator/purchases", middlewares.PermissionBlockerMiddleware("view_purchase_history"), controllers.AggregatorController{}.GetPurchases)
	//--------------------------------------------------------------------------------------------

	//----------------------------inventory-------------------------------------------------------------
	router.POST("aggregator/make_inventory_adjustment", middlewares.PermissionBlockerMiddleware("make_inventory_adjustment"), controllers.AggregatorController{}.MakeInventoryAdjustments)
	//--------------------------------------------------------------------------------------------------

	//---------------------------Aggregator ------------------------------------------------------
	router.POST("aggregator/add", middlewares.PermissionBlockerMiddleware("add_aggregator"), controllers.AggregatorController{}.InsertAggregator)
	router.GET("aggregators", middlewares.PermissionBlockerMiddleware("view_aggregator"), controllers.AggregatorController{}.GetAllAggregators)
	router.GET("aggregator/:id", middlewares.PermissionBlockerMiddleware("view_aggregator"), controllers.AggregatorController{}.GetAggregator)
	router.PUT("aggregator/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_aggregator"), controllers.AggregatorController{}.UpdateAggregatorStatus)
	router.DELETE("aggregator/delete/:id", middlewares.PermissionBlockerMiddleware("delete_aggregator"), controllers.AggregatorController{}.DeleteAggregator)
	router.PUT("aggregator/update", middlewares.PermissionBlockerMiddleware("edit_aggregator"), controllers.AggregatorController{}.UpdateAggregator)

	router.POST("aggregator/add/user", middlewares.PermissionBlockerMiddleware("add_user"), authController.AddAggregatorUser)
	router.PUT("aggregator/update/user", middlewares.PermissionBlockerMiddleware("edit_user"), authController.UpdateAggregatorUser)

	router.GET("aggregator/waste_types", middlewares.PermissionBlockerMiddleware("view_waste_type"), controllers.AggregatorController{}.GetWasteTypes)
	router.POST("aggregator/waste_types/create", middlewares.PermissionBlockerMiddleware("create_waste_type"), controllers.AggregatorController{}.SetWasteTypes)
	//-------------------------------------------------------------------------------------------

	//---------------------------Green champion ------------------------------------------------------
	router.POST("green_champion/add", middlewares.PermissionBlockerMiddleware("add_green_champion"), controllers.GreenChampionController{}.InsertGreenChampion)
	router.GET("green_champions", middlewares.PermissionBlockerMiddleware("view_green_champion"), controllers.GreenChampionController{}.GetAllGreenChampions)
	router.GET("green_champion/:id", middlewares.PermissionBlockerMiddleware("view_green_champion"), controllers.GreenChampionController{}.GetGreenChampion)
	router.PUT("green_champion/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_green_champion"), controllers.GreenChampionController{}.UpdateGreenChampionStatus)
	router.DELETE("green_champion/delete/:id", middlewares.PermissionBlockerMiddleware("delete_green_champion"), controllers.GreenChampionController{}.DeleteGreenChampion)
	router.PUT("green_champion/update", middlewares.PermissionBlockerMiddleware("edit_green_champion"), controllers.GreenChampionController{}.UpdateGreenChampion)
	//-------------------------------------------------------------------------------------------

	//---------------------------Roles ------------------------------------------------------
	// router.GET("roles", rolesController.GetRoles)
	// router.GET("role/:id", rolesController.GetRole)
	// router.POST("role/update", rolesController.UpdateRole)
	//-------------------------------------------------------------------------------------------

	//--------------------------- TTNM Organization ------------------------------------------------------
	router.GET("ttnm/profile/:id", middlewares.PermissionBlockerMiddleware("view_ttnm_profile"), ttnmOrganizationController.GetTTNMOrganizations)
	router.POST("ttnm/profile/update", middlewares.PermissionBlockerMiddleware("update_ttnm_profile"), ttnmOrganizationController.UpdateTtnmOrganizationProfile)
	//-------------------------------------------------------------------------------------------

	//--------------------------- Role Has Permissions ------------------------------------------------------
	// router.POST("permissions/assign", roleHasPermissionsController.AssignPermission)
	// router.POST("permissions/revoke", roleHasPermissionsController.RevokePermission)
	//roles and permissions -----------------------------------------------------------
	router.GET("roles", middlewares.PermissionBlockerMiddleware("view_roles"), controllers.RoleAndPermissionsController{}.GetRoles)
	router.GET("role/:id", middlewares.PermissionBlockerMiddleware("view_roles"), controllers.RoleAndPermissionsController{}.GetRole)
	router.DELETE("role/:id", middlewares.PermissionBlockerMiddleware("delete_roles"), controllers.RoleAndPermissionsController{}.DeleteRole)
	router.PUT("role/update", middlewares.PermissionBlockerMiddleware("edit_roles"), controllers.RoleAndPermissionsController{}.UpdateRole)
	router.POST("role/add", middlewares.PermissionBlockerMiddleware("add_roles"), controllers.RoleAndPermissionsController{}.AddRole)

	router.GET("permissions", middlewares.PermissionBlockerMiddleware("view_permissions"), controllers.RoleAndPermissionsController{}.GetAllPermissions)
	router.GET("permissions/:role_id", middlewares.PermissionBlockerMiddleware("view_permissions"), controllers.RoleAndPermissionsController{}.GetRolePermissions)
	router.GET("permissions/active_role_permissions/:role_id", middlewares.PermissionBlockerMiddleware("view_permissions"), controllers.RoleAndPermissionsController{}.GetActiveRolePermissions)

	router.PUT("assign_permissions_to_role", middlewares.PermissionBlockerMiddleware("assign_permissions_to_role"), controllers.RoleAndPermissionsController{}.AssignPermissionsToRole)
	router.PUT("remove_permissions_from_role", middlewares.PermissionBlockerMiddleware("remove_permissions_from_role"), controllers.RoleAndPermissionsController{}.RemovePermissionsFromRole)
	//---------------------------------------------------------------------------------
	//-------------------------------------------------------------------------------------------

	//--------------------------- wastegroups-----------------------------------------------------
	router.POST("settings/wastetypes/create", middlewares.PermissionBlockerMiddleware("create_waste_type"),wasteTypesController.InsertWasteGroup)
	router.PUT("settings/wastetypes/update", middlewares.PermissionBlockerMiddleware("update_waste_type"), wasteTypesController.UpdateWasteType)
	router.GET("settings/wastetypes/all",middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetAllWasteTypes)
	router.GET("settings/wastetypes/user",middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetUsersWasteGroups)
	router.GET("settings/wastetypes/wastegroup/:id", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetOneWasteGroup)
	//--------------------------------------------------------------------------------------------

	//--------------------------- Assign collectors to champions-----------------------------------------------------
	router.POST("assign_collectors_to_champions/assign", middlewares.PermissionBlockerMiddleware("assign_collector_to_champion"), championCollectorController.AssignAggregatorsToGreenChampionsParam)

	router.GET("get_collectors_for_green_champion/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetCollectorsForGreenChampion)
	router.GET("get_green_champions_for_collector/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetAllChampionsForACollector)
	//--------------------------------------------------------------------------------------------

	//--------------------------- Request Collections -----------------------------------------------------
	router.POST("request_collection", requestCollectionController.InsertNewCollectionRequestParams)
	router.POST("confirm_collection_request", requestCollectionController.ConfirmCollectionRequest)
	router.POST("cancel_collection_request", requestCollectionController.CancelCollectionRequest)
	router.POST("update_collection_request", requestCollectionController.UpdateCollectionRequest)
	router.GET("collections_weight_totals/:id", requestCollectionController.CollectionWeightTotals)
	//--------------------------------------------------------------------------------------------

	//--------------------------- Request Collections -----------------------------------------------------
	router.POST("collection_request_data", wasteItemsController.InsertWasteItem)
	router.GET("collection_request_latest/:id", requestCollectionController.GetLatestCollection)
	router.GET("collections_producer_waste_data/:id", requestCollectionController.GetWasteItemsProducerData)
	router.GET("collections_producer_stats/:id", requestCollectionController.GetCollectionStats)
	router.GET("collections_producer_complete/:id", requestCollectionController.GetAllProducerCompletedCollectionRequests)
	router.GET("collections_producer_pending/:id", requestCollectionController.GetAllProducerPendingCollectionRequests)
	//--------------------------------------------------------------------------------------------

	router.Run()
	//fmt.Println("Hello dabid")
}
