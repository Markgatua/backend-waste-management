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

	}

	apiGroup := router.Group("/api:BVaDN9hl")

	apiGroup.Use(middlewares.JwtAuthMiddleware())
	apiGroup.Use(middlewares.PermissionMiddleware())

	apiGroup.GET("/users", usersController.GetAllUsers)
	// apiGroup.POST("update/user",usersController.UpdateUSer)
	// apiGroup.GET("/users/roles",usersController.GetUsersWithRole)
	apiGroup.GET("/user/:id", usersController.GetUser)

	//---------------------------countries-------------------------------------------------------
	apiGroup.GET("countries", geoController.GetAllCountries)
	//-------------------------------------------------------------------------------------------

	//---------------------------organization----------------------------------------------------
	apiGroup.POST("organization/add", middlewares.PermissionBlockerMiddleware("add_organization"), organzationController.InsertOrganization)
	apiGroup.PUT("organization/update", middlewares.PermissionBlockerMiddleware("edit_organization"), organzationController.UpdateOrganization)
	apiGroup.GET("organizations", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetAllOrganizations)
	apiGroup.DELETE("organization/delete/:id", middlewares.PermissionBlockerMiddleware("delete_organization"), organzationController.DeleteOrganization)
	apiGroup.GET("organization/:id", middlewares.PermissionBlockerMiddleware("view_organizations"), organzationController.GetOrganization)
	//-------------------------------------------------------------------------------------------

	//---------------------------companies ------------------------------------------------------
	apiGroup.POST("companies/add", middlewares.PermissionBlockerMiddleware("add_company"), companiesController.InsertCompany)
	apiGroup.GET("companies", middlewares.PermissionBlockerMiddleware("view_companies"), companiesController.GetAllCompanies)
	apiGroup.GET("company/:id", middlewares.PermissionBlockerMiddleware("view_companies"), companiesController.GetCompany)
	apiGroup.POST("companies/status", middlewares.PermissionBlockerMiddleware("edit_company"), companiesController.UpdateCompanyStatus)
	apiGroup.DELETE("company/delete/:id", middlewares.PermissionBlockerMiddleware("delete_company"), companiesController.DeleteCompany)
	apiGroup.POST("companies/update", middlewares.PermissionBlockerMiddleware("edit_company"), companiesController.UpdateCompany)
	//-------------------------------------------------------------------------------------------

	//---------------------------Roles ------------------------------------------------------
	// apiGroup.GET("roles", rolesController.GetRoles)
	// apiGroup.GET("role/:id", rolesController.GetRole)
	// apiGroup.POST("role/update", rolesController.UpdateRole)
	//-------------------------------------------------------------------------------------------

	//--------------------------- TTNM Organization ------------------------------------------------------
	apiGroup.GET("ttnm/profile/:id", middlewares.PermissionBlockerMiddleware("view_ttnm_profile"), ttnmOrganizationController.GetTTNMOrganizations)
	apiGroup.POST("ttnm/profile/update", middlewares.PermissionBlockerMiddleware("update_ttnm_profile"), ttnmOrganizationController.UpdateTtnmOrganizationProfile)
	//-------------------------------------------------------------------------------------------

	//--------------------------- Role Has Permissions ------------------------------------------------------
	// apiGroup.POST("permissions/assign", roleHasPermissionsController.AssignPermission)
	// apiGroup.POST("permissions/revoke", roleHasPermissionsController.RevokePermission)
	//roles and permissions -----------------------------------------------------------
	apiGroup.GET("roles", middlewares.PermissionBlockerMiddleware("view_roles"), controllers.RoleAndPermissionsController{}.GetRoles)
	apiGroup.GET("role/:id", middlewares.PermissionBlockerMiddleware("view_roles"), controllers.RoleAndPermissionsController{}.GetRole)
	apiGroup.DELETE("role/:id", middlewares.PermissionBlockerMiddleware("delete_roles"), controllers.RoleAndPermissionsController{}.DeleteRole)
	apiGroup.PUT("role/update", middlewares.PermissionBlockerMiddleware("edit_roles"), controllers.RoleAndPermissionsController{}.UpdateRole)
	apiGroup.POST("role/add", middlewares.PermissionBlockerMiddleware("add_roles"), controllers.RoleAndPermissionsController{}.AddRole)

	apiGroup.GET("permissions", middlewares.PermissionBlockerMiddleware("view_permissions"), controllers.RoleAndPermissionsController{}.GetAllPermissions)
	apiGroup.GET("permissions/:role_id", middlewares.PermissionBlockerMiddleware("view_permissions"), controllers.RoleAndPermissionsController{}.GetRolePermissions)

	apiGroup.PUT("assign_permissions_to_role", middlewares.PermissionBlockerMiddleware("assign_permissions_to_role"), controllers.RoleAndPermissionsController{}.AssignPermissionsToRole)
	apiGroup.PUT("remove_permissions_from_role", middlewares.PermissionBlockerMiddleware("remove_permissions_from_role"), controllers.RoleAndPermissionsController{}.RemovePermissionsFromRole)
	//---------------------------------------------------------------------------------
	//-------------------------------------------------------------------------------------------

	//--------------------------- wastegroups-----------------------------------------------------
	apiGroup.POST("settings/wastegroups/create", middlewares.PermissionBlockerMiddleware("create_waste_type"), wasteGroupController.InsertWasteGroup)
	apiGroup.PUT("settings/wastegroups/update", middlewares.PermissionBlockerMiddleware("update_waste_type"), wasteGroupController.UpdateWasteGroup)
	apiGroup.GET("settings/wastegroups/all", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetAllWasteGroups)
	apiGroup.GET("settings/wastegroups/user", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetUsersWasteGroups)
	apiGroup.GET("settings/wastegroups/wastegroup/:id", middlewares.PermissionBlockerMiddleware("view_waste_type"), wasteGroupController.GetOneWasteGroup)
	//--------------------------------------------------------------------------------------------

	//--------------------------- Assign collectors to champions-----------------------------------------------------
	apiGroup.POST("assign_collectors_to_champions/assign", middlewares.PermissionBlockerMiddleware("assign_collector_to_champion"), championCollectorController.AssignChampionToCollector)
	apiGroup.GET("assign_collectors_to_champions/get_champion_collector/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetTheCollectorForAChampion)
	apiGroup.GET("assign_collectors_to_champions/get_champions_for_a_collector/:id", middlewares.PermissionBlockerMiddleware("view_champion_collector"), championCollectorController.GetAllChampionsForACollector)
	apiGroup.POST("assign_collectors_to_champions/update", middlewares.PermissionBlockerMiddleware("assign_collector_to_champion"), championCollectorController.UpdateChampionCollector) //not in use
	// apiGroup.POST("settings/wastegroups/update", wasteGroupController.UpdateWasteGroup)
	// apiGroup.GET("settings/wastegroups/all", wasteGroupController.GetAllWasteGroups)
	// apiGroup.GET("settings/wastegroups/wastegroup/:id", wasteGroupController.GetOneWasteGroup)
	//--------------------------------------------------------------------------------------------


	router.Run()
	//fmt.Println("Hello dabid")
}
