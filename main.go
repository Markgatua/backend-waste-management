package main

import (
	"fmt"
	_ "fmt"
	"os"
	"ttnmwastemanagementsystem/controllers"
	"ttnmwastemanagementsystem/database/seeder"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
)

func main() {
	var instructions string = `
	   Instructions
	   ------------
	   1) To run seeder run "go run main.go seeder"
	   2) To run master table migration run "go run main.go master_table_migration"
	   3) To run the school system run "go run main.go scms"
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
	gen.LoadRepo()

	usersController := controllers.UsersController{}
	companiesController := controllers.CompaniesController{}
	wasteGroupController := controllers.WasteGroupsController{}
	companyRegionsController := controllers.CompanyRegionsController{}
	companyBranchesController := controllers.CompanyBranchesController{}
	apiGroup := router.Group("/api:BVaDN9hl")

	apiGroup.GET("/users", usersController.GetAllUsers)
	// apiGroup.POST("update/user",usersController.UpdateUSer)
	// apiGroup.GET("/users/roles",usersController.GetUsersWithRole)
	apiGroup.GET("/user/:id", usersController.GetUser)
	apiGroup.POST("companies/addcompany", companiesController.InsertCompany)
	apiGroup.GET("companies/allcompanies", companiesController.GetAllCompanies)
	apiGroup.GET("companies/company/:id", companiesController.GetCompany)
	apiGroup.POST("companies/status", companiesController.UpdateCompanyStatus)
	apiGroup.POST("companies/update", companiesController.UpdateCompany)
	apiGroup.POST("companies/regions/create", companyRegionsController.InsertCompanyRegion)
	apiGroup.GET("companies/company/regions/:company_id", companyRegionsController.GetAllCompanyRegions)
	apiGroup.GET("companies/regions", companyRegionsController.GetAllCompaniesRegions)
	apiGroup.GET("companies/company/region/:id", companyRegionsController.GetOneCompanyRegion)
	apiGroup.POST("companies/regions/update", companyRegionsController.UpdateCompanyRegionData)
	apiGroup.POST("companies/branch/create", companyBranchesController.InsertCompanyBranch)
	apiGroup.POST("companies/branch/status", companyBranchesController.UpdateCompanyBranchStatus)
	apiGroup.POST("settings/wastegroups/create", wasteGroupController.InsertWasteGroup)
	apiGroup.POST("settings/wastegroups/update", wasteGroupController.UpdateWasteGroup)
	apiGroup.GET("settings/wastegroups/all", wasteGroupController.GetAllWasteGroups)
	apiGroup.GET("settings/wastegroups/wastegroup/:id", wasteGroupController.GetOneWasteGroup)
	router.Run()
	//fmt.Println("Hello dabid")
}
