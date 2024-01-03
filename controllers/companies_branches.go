package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
	_ "github.com/guregu/null"
)


type CompanyBranchesController struct{}

type InsertCompanyBranchParams struct {
	Branch          string `json:"name"  binding:"required"`
	CompanyID       int32 `json:"company_id"  binding:"required"`
	RegionID        int32 `json:"region_id"`
	BranchLocation  string `json:"location"  binding:"required"`
}

type UpdateCompanyBranchDataParams struct {
	ID     		  int `json:"id"  binding:"required"`
	RegionID        int32 `json:"region_id"`
	Branch          string `json:"name"  binding:"required"`
	BranchLocation  string `json:"location"  binding:"required"`
}

type UpdateCompanyBranchStatusParams struct {
	ID     	int `json:"id"  binding:"required"`
	IsActive 	string `json:"status"  binding:"required"`
}

func (companyBranchesController CompanyBranchesController) InsertCompanyBranch(context *gin.Context) {
	var params InsertCompanyBranchParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	companyBranch, insertError := gen.REPO.InsertCompanyBranch(context, gen.InsertCompanyBranchParams{
		CompanyID:   	int32(params.CompanyID),
		RegionID:    	sql.NullInt32{Int32: int32(params.RegionID), Valid: params.RegionID != 0},
		Branch:      	params.Branch,
		BranchLocation: params.BranchLocation,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Company Branch",
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Company Branch",
		"company branch": companyBranch, // Include the company details in the response
	})
}


func (companyBranchesController CompanyBranchesController) UpdateCompanyBranchStatus(context *gin.Context) {
	var params UpdateCompanyBranchStatusParams
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
	updateError := gen.REPO.UpdateCompanyBranchStatus(context, gen.UpdateCompanyBranchStatusParams{
		IsActive: status,
		ID:     int32(params.ID),
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
		"message": "Successfully updated Company Branch Status",
		"status":  status, // Use the variable for the parsed status
	})
}

