package main

import (
	_ "fmt"
	"ttnmwastemanagementsystem/controllers"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)
func main()  {
	router := gin.Default()
	gen.LoadRepo()

	usersController:=controllers.UsersController{}
	companiesController:=controllers.CompaniesController{}
	wasteGroupController:=controllers.WasteGroupsController{}
	companyRegionsController:=controllers.CompanyRegionsController{}
	apiGroup := router.Group("/api:BVaDN9hl")

	apiGroup.GET("/users",usersController.GetAllUsers)
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
	apiGroup.POST("settings/wastegroups/create", wasteGroupController.InsertWasteGroup)
	apiGroup.POST("settings/wastegroups/update", wasteGroupController.UpdateWasteGroup)
	apiGroup.GET("settings/wastegroups/all", wasteGroupController.GetAllWasteGroups)
	apiGroup.GET("settings/wastegroups/wastegroup/:id", wasteGroupController.GetOneWasteGroup)
	router.Run()
	//fmt.Println("Hello dabid")
}