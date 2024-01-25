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
	//geoController := controllers.GeoController{}
	championCollectorController := controllers.ChampionCollectorController{}
	// rolesController := controllers.RolesController{}
	// rolespermissions := controllers.RoleAndPermissionsController{}
	ttnmOrganizationController := controllers.TtnmOrganizationController{}

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

	}

	// apiGroup := router.Group("/api:BVaDN9hl")
	router.GET("/uploads/:file", controllers.FileController{}.GetFile)
	router.GET("/flag/:file", controllers.FileController{}.GetFlag)

	//presets ----------------------------------------------------------|
	router.GET("/presets", controllers.PresetController{}.GetPresetValue)
	//------------------------------------------------------------------|

	router.Use(middlewares.JwtAuthMiddleware())
	router.Use(middlewares.PermissionMiddleware())

	//---------------------------    Files ------------------------------------------------------
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload_files", middlewares.PermissionBlockerMiddleware("upload_files"), controllers.FileController{}.UploadFiles)
	//-------------------------------------------------------------------------------------------

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

	//---------------------------countries-------------------------------------------------------
	//router.GET("countries", geoController.GetAllCountries)
	//-------------------------------------------------------------------------------------------

	//---------------------------organization----------------------------------------------------
	router.POST("organization/add", middlewares.PermissionBlockerMiddleware("add_organization"), organzationController.InsertOrganization)
	router.PUT("organization/update", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.UpdateOrganization)
	router.GET("organizations", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetAllOrganizations)
	router.PUT("organization/set_active_inactive_status", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.SetActiveInActiveStatus)

	router.DELETE("organization/delete/:id", middlewares.PermissionBlockerMiddleware("delete_organization"), organzationController.DeleteOrganization)
	router.GET("organization/:id", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetOrganization)
	//-------------------------------------------------------------------------------------------

	//---------------------------Aggregator ------------------------------------------------------
	router.POST("aggregator/add", middlewares.PermissionBlockerMiddleware("add_aggregator"), controllers.AggregatorController{}.InsertCompany)
	router.GET("aggregators", middlewares.PermissionBlockerMiddleware("view_aggregator"), controllers.AggregatorController{}.GetAllCompanies)
	router.GET("aggregator/:id", middlewares.PermissionBlockerMiddleware("view_aggregator"), controllers.AggregatorController{}.GetCompany)
	router.POST("aggregator/status", middlewares.PermissionBlockerMiddleware("edit_aggregator"), controllers.AggregatorController{}.UpdateCompanyStatus)
	router.DELETE("aggregator/delete/:id", middlewares.PermissionBlockerMiddleware("delete_aggregator"), controllers.AggregatorController{}.DeleteCompany)
	router.POST("aggregator/update", middlewares.PermissionBlockerMiddleware("edit_aggregator"), controllers.AggregatorController{}.UpdateCompany)
	//-------------------------------------------------------------------------------------------

	//---------------------------Green Champion -------------------------------------------------
	router.POST("green_champion/add", middlewares.PermissionBlockerMiddleware("add_green_champion"), controllers.AggregatorController{}.InsertCompany)
	router.GET("green_champion", middlewares.PermissionBlockerMiddleware("view_green_champion"), controllers.AggregatorController{}.GetAllCompanies)
	router.GET("green_champion/:id", middlewares.PermissionBlockerMiddleware("view_green_champion"), controllers.AggregatorController{}.GetCompany)
	router.POST("green_champion/status", middlewares.PermissionBlockerMiddleware("edit_green_champion"), controllers.AggregatorController{}.UpdateCompanyStatus)
	router.DELETE("green_champion/delete/:id", middlewares.PermissionBlockerMiddleware("delete_green_champion"), controllers.AggregatorController{}.DeleteCompany)
	router.POST("green_champion/update", middlewares.PermissionBlockerMiddleware("edit_green_champion"), controllers.AggregatorController{}.UpdateCompany)
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
	router.POST("settings/wastetypes/create", middlewares.PermissionBlockerMiddleware("create_waste_type"), wasteTypesController.InsertWasteGroup)
	router.PUT("settings/wastetypes/update", middlewares.PermissionBlockerMiddleware("update_waste_type"), wasteTypesController.UpdateWasteType)
	router.GET("settings/wastetypes/all", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetAllWasteTypes)
	router.GET("settings/wastetypes/user", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetUsersWasteGroups)
	router.GET("settings/wastetypes/wastegroup/:id", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteTypesController.GetOneWasteGroup)
	//--------------------------------------------------------------------------------------------

	//--------------------------- Assign collectors to champions-----------------------------------------------------
	router.POST("assign_collectors_to_champions/assign", middlewares.PermissionBlockerMiddleware("assign_collector_to_champion"), championCollectorController.AssignChampionToCollector)
	router.GET("assign_collectors_to_champions/get_champion_collector/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetTheCollectorForAChampion)
	router.GET("assign_collectors_to_champions/get_champions_for_a_collector/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetAllChampionsForACollector)
	router.POST("assign_collectors_to_champions/update", middlewares.PermissionBlockerMiddleware("assign_collector_to_champion"), championCollectorController.UpdateChampionCollector) //not in use
	// apiGroup.POST("settings/wastegroups/update", wasteGroupController.UpdateWasteGroup)
	// apiGroup.GET("settings/wastegroups/all", wasteGroupController.GetAllWasteGroups)
	// apiGroup.GET("settings/wastegroups/wastegroup/:id", wasteGroupController.GetOneWasteGroup)
	//--------------------------------------------------------------------------------------------

	router.Run()
	//fmt.Println("Hello dabid")
}
