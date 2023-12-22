package controllers

import (
	"net/http"
	"strconv"
	_ "strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)


type CompanyRegionsController struct{}

type InsertCompanyRegionParams struct {
	Region         string `json:"name"  binding:"required"`
	CompanyID      int32 `json:"company_id"  binding:"required"`
}

type UpdateCompanyRegionDataParams struct {
	ID     		  int `json:"id"  binding:"required"`
	Region 		  string `json:"name"  binding:"required"`
}


func (companyRegionsController CompanyRegionsController) InsertCompanyRegion(context *gin.Context) {
	var params InsertCompanyRegionParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}



	companyRegion, insertError := gen.REPO.InsertCompanyRegion(context, gen.InsertCompanyRegionParams{
		Region:		 params.Region,
		CompanyID:   int32(params.CompanyID),	
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Company Region",
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Company Region",
		"company": companyRegion, // Include the company details in the response
	})
}

func (c *CompanyRegionsController) GetAllCompanyRegions(context *gin.Context) {
    // Extract company_id from the request URL
    CompanyID := context.Param("company_id")

    // Convert company_id to uint
    CompanyID_, err := strconv.ParseUint(CompanyID, 10, 32)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "error":   true,
            "message": "Invalid company_id",
        })
        return
    }

    // Convert uint to int32
    CompanyIDInt32 := int32(CompanyID_)

    // Call the repository method with the int32 CompanyID
    companyRegions, err := gen.REPO.GetAllCompanyRegions(context, CompanyIDInt32)
    if err != nil {
        context.JSON(http.StatusUnprocessableEntity, gin.H{
            "error":   true,
            "message": err.Error(),
        })
        return
    }

    context.JSON(http.StatusOK, gin.H{
        "error":        false,
        "Company Regions": companyRegions,
    })
}


func(companyRegionsController CompanyRegionsController) GetOneCompanyRegion(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	companyRegion, err := gen.REPO.GetOneCompanyRegion(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"Company Regions": companyRegion,
	})
}

func(companyRegionsController  CompanyRegionsController) GetAllCompaniesRegions(context *gin.Context){
	companies_regions, err := gen.REPO.GetAllCompaniesRegions(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"Companies Regions":companies_regions,
	})
}


func (companyRegionsController CompanyRegionsController) UpdateCompanyRegionData(context *gin.Context) {
	var params UpdateCompanyRegionDataParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update Waste Group
	updateError := gen.REPO.UpdateCompanyRegionData(context, gen.UpdateCompanyRegionDataParams{
		Region: params.Region,
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
		"message": "Successfully updated Company Region",
	})
}
