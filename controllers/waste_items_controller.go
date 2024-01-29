package controllers

import (
	"fmt"
	"net/http"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type WasteItemsController struct{}

type InsertWasteItemParams struct {
	CollectionRequestID int32 `json:"collection_request_id"`
	WasteTypeID         int32 `json:"waste_type_id"`
	Weight              string        `json:"weight"`
}

func (wasteItemsController WasteItemsController) InsertWasteItem(context *gin.Context) {
	var params InsertWasteItemParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	wasteType, insertError := gen.REPO.InsertWasteItem(context, gen.InsertWasteItemParams{
		CollectionRequestID: params.CollectionRequestID,
		WasteTypeID: params.WasteTypeID,
		Weight: params.Weight,
	})

	fmt.Println(insertError)

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add skuS",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Successfully Added Sku",
		"waste_type": wasteType, 
	})

}
