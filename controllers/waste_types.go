package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)


type WasteTypesController struct{}

type InsertWasteGroupParams struct {
	Name      string `json:"name"  binding:"required"`
	Category      string `json:"category"`
}

type UpdateWasteGroupParams struct {
	ID     		  int `json:"id"  binding:"required"`
	Name 		  string `json:"name"  binding:"required"`
	Category      string `json:"category"`
	IsActive	  *bool   `json:"is_active"`
}


func (wasteGroupsController WasteTypesController) InsertWasteGroup(context *gin.Context) {
	var params InsertWasteGroupParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	WasteGroup, insertError := gen.REPO.InsertWasteType(context, gen.InsertWasteTypeParams{
		Name:		 params.Name,
		Category:    params.Category,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Waste Type",
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Waste Type",
		"waste_type": WasteGroup, // Include the company details in the response
	})

}

func(wasteGroupsController  WasteTypesController) GetAllWasteTypes(context *gin.Context){
	wasteGroups, err := gen.REPO.GetAllWasteTypes(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"waste_types":wasteGroups,
	})
}

func(wasteGroupsController  WasteTypesController) GetUsersWasteGroups(context *gin.Context){
	wasteGroups, err := gen.REPO.GetUsersWasteType(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"waste_types":wasteGroups,
	})
}

func(wasteGroupController WasteTypesController) GetOneWasteGroup(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	wasteGroup, err := gen.REPO.GetOneWasteType(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"waste_type": wasteGroup,
	})
}


func (wasteGroupController WasteTypesController) UpdateWasteGroup(context *gin.Context) {
	var params UpdateWasteGroupParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	
	// Update Waste Group
	updateError := gen.REPO.UpdateWasteType(context, gen.UpdateWasteTypeParams{
		Category: params.Category,
		Name: params.Name,
		ID:     int32(params.ID),
		IsActive: *params.IsActive,
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
		"message": "Successfully updated Waste Group",
	})
}
