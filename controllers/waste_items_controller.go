package controllers

import (
	_"fmt"
	"net/http"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type WasteItemsController struct{}

type InsertWasteItemParams struct {
	CollectionRequestID int32 `json:"collection_request_id"`
	Waste []Waste `json:"waste"`
}
type Waste struct {
	WasteTypeID int64  `json:"waste_type_id"`
	Weight      string `json:"weight"`
}

func (wasteItemsController WasteItemsController) InsertWasteItem(context *gin.Context) {
	var params InsertWasteItemParams

	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Iterate through waste items and insert each one
	for _, wasteItem := range params.Waste {
		collectedWaste,insertError := gen.REPO.InsertWasteItem(context, gen.InsertWasteItemParams{
			CollectionRequestID: params.CollectionRequestID,
			WasteTypeID:         int32(wasteItem.WasteTypeID),
			Weight:              wasteItem.Weight,
		})

		if insertError != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Failed to insert waste items",
				"Collected Waste": collectedWaste,
			})
			return
		}
	}
}
