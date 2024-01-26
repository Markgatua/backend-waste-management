package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type AggregatorController struct{}

type CreateAggregatorParams struct {
	CountyID         int32  `json:"county_id"  binding:"required"`
	PhysicalPosition string `json:"physical_position" binding:"required"`
	Name             string `json:"name"  binding:"required"`
	Location         string `json:"location"  binding:"required"`
	Email            string `json:"email" binding:"required"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Password         string `json:"password" binding:"required"`
	OrganizationID   int32  `json:"organization_id"`
	LogoPath         string `json:"logo_path"`

	IsActive *bool  `json:"is_active"  binding:"required"`
	Region   string `json:"region"  binding:"required"`
}

type UpdateAggregatorDataParams struct {
	CountyID         int32  `json:"county_id"  binding:"required"`
	PhysicalPosition string `json:"physical_position" binding:"required"`
	Name             string `json:"name"  binding:"required"`
	Companytype      int32  `json:"company_type"  binding:"required"`
	Location         string `json:"location"  binding:"required"`
	OrganizationID   int32  `json:"organization_id"  binding:"required"`
	IsActive         *bool  `json:"is_active"  binding:"required"`
	ID               int64  `json:"id"  binding:"required"`
	LogoPath         string `json:"logo_path"`
	Region           string `json:"region"  binding:"required"`
	UserID           int32  `json:"user_id" binding:"required"`
	Email            string `json:"email" binding:"required"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Password         string `json:"password"`
}

type UpdateAggregatorStatusParams struct {
	ID       int   `json:"id"  binding:"required"`
	IsActive *bool `json:"status"  binding:"required"`
}

func (controller AggregatorController) InsertAggregator(context *gin.Context) {
	var params CreateAggregatorParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := GetEmailUser(params.Email)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	context.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error":   true,
	// 		"message": "Error adding organization",
	// 	})
	// 	return
	// }
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email already exists`",
		})
		return
	}

	count, err := gen.REPO.GetDuplicateCompanies(context, gen.GetDuplicateCompaniesParams{
		Name:           strings.ToLower(params.Name),
		OrganizationID: sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another company has the same name withing the organization",
		})
		return
	}

	company, insertError := gen.REPO.InsertCompany(context, gen.InsertCompanyParams{
		CountyID:         sql.NullInt32{Int32: params.CountyID, Valid: params.CountyID != 0},
		PhysicalPosition: params.PhysicalPosition,
		Name:             params.Name,
		Location:         null.StringFrom(params.Location).NullString,
		IsActive:         *params.IsActive,
		OrganizationID:   sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
		Region:           null.StringFrom(params.Region).NullString,
		CompanyType:      2,//params.Companytype,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add aggregator",
		})
		return
	}

	_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,provider,role_id,user_organization_id,user_type,created_at,updated_at,password,is_organization_super_admin,user_company_id) VALUES (:email,:first_name,:last_name,:provider,:role_id,:user_organization_id,:user_type,:created_at,:updated_at,:password,:is_organization_super_admin,:user_company_id)`,
		map[string]interface{}{
			"email":                       params.Email,
			"first_name":                  params.FirstName,
			"last_name":                   params.LastName,
			"provider":                    "email",
			"role_id":                     3,
			"user_organization_id":        params.OrganizationID,
			"is_organization_super_admin": true,
			"user_type":                   9,
			"user_company_id":             company.ID,
			"password":                    helpers.Functions{}.HashPassword(params.Password),
			"created_at":                  time.Now(),
			"updated_at":                  time.Now(),
		})

	if params.LogoPath != "" {
		UploadController{}.SaveToUploadsTable(params.LogoPath, "companies", company.ID)
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Aggregator",
		"company": company, // Include the company details in the response
	})
}

func (controller AggregatorController) GetAllAggregators(context *gin.Context) {
	companies, err := gen.REPO.GetAllAggregators(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": companies,
	})
}

func (controller AggregatorController) GetAggregator(context *gin.Context) {
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

func (c AggregatorController) DeleteAggregator(context *gin.Context) {
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

func (aggregatorController AggregatorController) UpdateAggregatorStatus(context *gin.Context) {
	var params UpdateAggregatorStatusParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateCompanyStatus(context, gen.UpdateCompanyStatusParams{
		IsActive: *params.IsActive,
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
		"message": "Successfully updated Aggregator",
		"status":  params.IsActive, // Use the variable for the parsed status
	})
}

func (aggregatorController AggregatorController) UpdateAggregator(context *gin.Context) {
	var params UpdateAggregatorDataParams
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
		OrganizationID: sql.NullInt32{Int32: params.OrganizationID, Valid: true},
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
		CountyID:         sql.NullInt32{Int32: params.CountyID, Valid: params.CountyID != 0},
		PhysicalPosition: params.PhysicalPosition,
		IsActive:         *params.IsActive,
		Name:             params.Name,
		Location:         null.StringFrom(params.Location).NullString,
		Region:           null.StringFrom(params.Region).NullString,
		OrganizationID:   sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
		CompanyType:      params.Companytype,
		ID:               int32(params.ID),
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
