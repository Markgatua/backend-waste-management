package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type CompaniesController struct{}


type CreateCompanyParams struct {
	Name      string `json:"name"  binding:"required"`
	Companytype int32 `json:"companytype"  binding:"required"`
	Logo      string `json:"logo"`
	Location  string `json:"location"  binding:"required"`
	HasRegions  string   `json:"has_regions"  binding:"required"`
	HasBranches  string   `json:"has_branches"  binding:"required"`

}

type UpdateCompanyDataParams struct {
	ID        int    `json:"id"`
	Name      string `json:"name"  binding:"required"`
	Logo      string `json:"logo"`
	Location  string `json:"location"  binding:"required"`
	HasRegions  string   `json:"has_regions"  binding:"required"`
	HasBranches  string   `json:"has_branches"  binding:"required"`
}

type UpdateCompanyStatusParams struct {
	ID     	int `json:"id"  binding:"required"`
	IsActive 	string `json:"status"  binding:"required"`
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


	// Convert HasRegions from string to bool
	hasRegions, err := strconv.ParseBool(params.HasRegions)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid value for 'has_regions'",
		})
		return
	}

	// Convert HasBranches from string to bool
	hasBranches, err := strconv.ParseBool(params.HasBranches)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid value for 'has_branches'",
		})
		return
	}


	company, insertError := gen.REPO.InsertCompany(context, gen.InsertCompanyParams{
		Name:		 params.Name,
		Location:    null.StringFrom(params.Location).NullString,
		Companytype: params.Companytype,
		HasRegions:  hasRegions,
		HasBranches: hasBranches,
		Logo:        null.StringFrom(params.Logo).NullString,
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


func(companiesController  CompaniesController) GetAllCompanies(context *gin.Context){
	companies, err := gen.REPO.GetAllCompanies(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"companies":companies,
	})
}

func(companiesController CompaniesController) GetCompany(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	company, err := gen.REPO.GetCompany(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"company": company,
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

	// Convert HasRegions from string to bool
	hasRegions, err := strconv.ParseBool(params.HasRegions)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid value for 'has_regions'",
		})
		return
	}

	// Convert HasBranches from string to bool
	hasBranches, err := strconv.ParseBool(params.HasBranches)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid value for 'has_branches'",
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateCompanyData(context, gen.UpdateCompanyDataParams{
		HasBranches: hasBranches,
		HasRegions: hasRegions,
		Name: params.Name,
		Location: null.StringFrom(params.Location).NullString,
		Logo: null.StringFrom(params.Logo).NullString,
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
		"message": "Successfully updated Company Data",
	})
}