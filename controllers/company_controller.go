package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type CompaniesController struct{}

type CreateCompanyParams struct {
	Name           string `json:"name"  binding:"required"`
	Companytype    int32  `json:"company_type"  binding:"required"`
	Location       string `json:"location"  binding:"required"`
	OrganizationID int32  `json:"organization_id"  binding:"required"`
	IsActive       bool   `json:"is_active"  binding:"required"`
	Region         string `json:"region"  binding:"required"`
}

type UpdateCompanyDataParams struct {
	Name           string `json:"name"  binding:"required"`
	Companytype    int32  `json:"company_type"  binding:"required"`
	Location       string `json:"location"  binding:"required"`
	OrganizationID int32  `json:"organization_id"  binding:"required"`
	IsActive       bool   `json:"is_active"  binding:"required"`
	ID             int64  `json:"id"  binding:"required"`
	Region         string `json:"region"  binding:"required"`
}

type UpdateCompanyStatusParams struct {
	ID       int    `json:"id"  binding:"required"`
	IsActive string `json:"status"  binding:"required"`
}

func (companiesController CompaniesController) InsertCompany(context *gin.Context) {
	var params CreateCompanyParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	count, err := gen.REPO.GetDuplicateCompanies(context, gen.GetDuplicateCompaniesParams{
		Name:           strings.ToLower(params.Name),
		OrganizationID: params.OrganizationID,
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another company has the same name withing the organization",
		})
		return
	}

	company, insertError := gen.REPO.InsertCompany(context, gen.InsertCompanyParams{
		Name:        params.Name,
		Location:    null.StringFrom(params.Location).NullString,
		IsActive:    params.IsActive,
		Region:      null.StringFrom(params.Region).NullString,
		CompanyType: params.Companytype,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Company",
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Company",
		"company": company, // Include the company details in the response
	})
}

func (companiesController CompaniesController) GetAllCompanies(context *gin.Context) {
	companies, err := gen.REPO.GetAllCompanies(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":     false,
		"companies": companies,
	})
}

func (companiesController CompaniesController) GetCompany(context *gin.Context) {
	id := context.Param("id")

	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	company, err := gen.REPO.GetCompany(context, int32(id_))

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"company": company,
	})
}

func (c OrgnizationController) DeleteCompany(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	err := gen.REPO.DeleteCompany(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Company successfully deleted",
	})
}

func (companiesController CompaniesController) UpdateCompanyStatus(context *gin.Context) {
	var params UpdateCompanyStatusParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Convert Status from string to bool
	status, err := strconv.ParseBool(params.IsActive)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid value for 'status'",
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateCompanyStatus(context, gen.UpdateCompanyStatusParams{
		IsActive: status,
		ID:       int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated Company Status",
		"status":  status, // Use the variable for the parsed status
	})
}

func (companiesController CompaniesController) UpdateCompany(context *gin.Context) {
	var params UpdateCompanyDataParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	count, err := gen.REPO.GetDuplicateCompaniesWithoutID(context, gen.GetDuplicateCompaniesWithoutIDParams{
		Name:           strings.ToLower(params.Name),
		ID:             int32(params.ID),
		OrganizationID: params.OrganizationID,
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another company has the same name withing the organization",
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateCompany(context, gen.UpdateCompanyParams{
		IsActive:       params.IsActive,
		Name:           params.Name,
		Location:       null.StringFrom(params.Location).NullString,
		Region:         null.StringFrom(params.Region).NullString,
		OrganizationID: params.OrganizationID,
		CompanyType:    params.Companytype,
		ID:             int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated Company Data",
	})
}
