package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)


type WasteGroupsController struct{}

type InsertWasteGroupParams struct {
	Name      string `json:"name"  binding:"required"`
	Category      string `json:"category"  binding:"required"`
}

type UpdateWasteGroupParams struct {
	ID     		  int `json:"id"  binding:"required"`
	Name 		  string `json:"name"  binding:"required"`
	Category      string `json:"category"  binding:"required"`
}


func (wasteGroupsController WasteGroupsController) InsertWasteGroup(context *gin.Context) {
	var params InsertWasteGroupParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	WasteGroup, insertError := gen.REPO.InsertWasteGroup(context, gen.InsertWasteGroupParams{
		Name:		 params.Name,
		Category:    params.Category,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Waste Group",
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Waste Group",
		"company": WasteGroup, // Include the company details in the response
	})

}

func(wasteGroupsController  WasteGroupsController) GetAllWasteGroups(context *gin.Context){
	wasteGroups, err := gen.REPO.GetAllWasteGroups(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"Waste Groups":wasteGroups,
	})
}

func(wasteGroupController WasteGroupsController) GetOneWasteGroup(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	wasteGroup, err := gen.REPO.GetOneWasteGroup(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"Waste Group": wasteGroup,
	})
}


func (wasteGroupController WasteGroupsController) UpdateWasteGroup(context *gin.Context) {
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
	updateError := gen.REPO.UpdateWasteGroup(context, gen.UpdateWasteGroupParams{
		Category: params.Category,
		Name: params.Name,
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
		"message": "Successfully updated Waste Group",
	})
}
