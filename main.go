package main

import (
	"fmt"
	_ "fmt"
	"os"
	_"time"
	_ "time"
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/controllers"
	"ttnmwastemanagementsystem/database/seeder"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/middlewares"

	cors "github.com/rs/cors/wrapper/gin"
	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
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
	companiesController := controllers.CompaniesController{}
	wasteGroupController := controllers.WasteGroupsController{}
	organzationController := controllers.OrgnizationController{}
	geoController := controllers.GeoController{}
	championCollectorController := controllers.ChampionCollectorController{}
	// rolesController := controllers.RolesController{}
	// rolespermissions := controllers.RoleAndPermissionsController{}
	ttnmOrganizationController := controllers.TtnmOrganizationController{}
	requestCollectionController := controllers.RequestCollectionController{}

	router.LoadHTMLGlob("templates/**/*")

	router.Use(cors.Default())


	auth := router.Group("/auth")
	{
		authController := controllers.AuthController{}
		auth.GET("/challenge/password_reset_success/email/web", authController.PassWordResetSuccess)
		auth.POST("/challenge/submit_new_password/email/web", authController.SubmitNewPassword)
		auth.GET("/challenge/enter_new_password/web", authController.EnterNewPassword)
		auth.POST("/reset_password/email/web", authController.ResetPassword)

		auth.POST("/register/email", authController.RegisterUserEmail)
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

	router.Use(middlewares.JwtAuthMiddleware())
	router.Use(middlewares.PermissionMiddleware())

	router.GET("/users", usersController.GetAllUsers)
	// router.POST("update/user",usersController.UpdateUSer)
	// router.GET("/users/roles",usersController.GetUsersWithRole)
	router.GET("/user/:id", usersController.GetUser)
	router.GET("/company/users/:id", usersController.GetCompanyUsers)

	//---------------------------countries-------------------------------------------------------
	router.GET("countries", geoController.GetAllCountries)
	//-------------------------------------------------------------------------------------------

	//---------------------------organization----------------------------------------------------
	router.POST("organization/add", middlewares.PermissionBlockerMiddleware("add_organization"), organzationController.InsertOrganization)
	router.PUT("organization/update", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.UpdateOrganization)
	router.GET("organizations", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetAllOrganizations)
	router.DELETE("organization/delete/:id", middlewares.PermissionBlockerMiddleware("delete_organization"), organzationController.DeleteOrganization)
	router.GET("organization/:id", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetOrganization)
	//-------------------------------------------------------------------------------------------

	//---------------------------companies ------------------------------------------------------
	router.POST("companies/add", middlewares.PermissionBlockerMiddleware("add_company"), companiesController.InsertCompany)
	router.GET("companies", middlewares.PermissionBlockerMiddleware("view_companies"), companiesController.GetAllCompanies)
	router.GET("company/:id", middlewares.PermissionBlockerMiddleware("view_companies"), companiesController.GetCompany)
	router.POST("companies/status", middlewares.PermissionBlockerMiddleware("edit_company"), companiesController.UpdateCompanyStatus)
	router.DELETE("company/delete/:id", middlewares.PermissionBlockerMiddleware("delete_company"), companiesController.DeleteCompany)
	router.POST("companies/update", middlewares.PermissionBlockerMiddleware("edit_company"), companiesController.UpdateCompany)
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

	router.PUT("assign_permissions_to_role", middlewares.PermissionBlockerMiddleware("assign_permissions_to_role"), controllers.RoleAndPermissionsController{}.AssignPermissionsToRole)
	router.PUT("remove_permissions_from_role", middlewares.PermissionBlockerMiddleware("remove_permissions_from_role"), controllers.RoleAndPermissionsController{}.RemovePermissionsFromRole)
	//---------------------------------------------------------------------------------
	//-------------------------------------------------------------------------------------------

	//--------------------------- wastegroups-----------------------------------------------------
	router.POST("settings/wastegroups/create", middlewares.PermissionBlockerMiddleware("create_waste_type"), wasteGroupController.InsertWasteGroup)
	router.PUT("settings/wastegroups/update", middlewares.PermissionBlockerMiddleware("update_waste_type"), wasteGroupController.UpdateWasteGroup)
	router.GET("settings/wastegroups/all", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetAllWasteGroups)
	router.GET("settings/wastegroups/user", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetUsersWasteGroups)
	router.GET("settings/wastegroups/wastegroup/:id", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetOneWasteGroup)
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

	//--------------------------- Request Collections -----------------------------------------------------
	router.POST("request_collection", middlewares.PermissionBlockerMiddleware("request_collection"), requestCollectionController.InsertNewCollectionRequestParams)
	router.POST("confirm_collection_request", requestCollectionController.ConfirmCollectionRequest)
	router.POST("cancel_collection_request", requestCollectionController.CancelCollectionRequest)
	router.POST("update_collection_request", requestCollectionController.UpdateCollectionRequest)
	//--------------------------------------------------------------------------------------------


	router.Run()
	//fmt.Println("Hello dabid")
}
