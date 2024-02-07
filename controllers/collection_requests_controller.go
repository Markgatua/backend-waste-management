package controllers

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"net/http"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"

	"github.com/gin-gonic/gin"
)

type CollectionRequestsController struct{}

type InsertWasteItemParams struct {
	CollectionRequestID int32   `json:"collection_request_id"`
	Waste               []Waste `json:"waste"`
}
type Waste struct {
	WasteTypeID int64   `json:"waste_type_id"`
	Weight      float64 `json:"weight"`
}

func (controller CollectionRequestsController) InsertWasteItems(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	var params InsertWasteItemParams

	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	errorInsertingWasteItem := false
	// Iterate through waste items and insert each one
	for _, wasteItem := range params.Waste {
		_, insertError := gen.REPO.InsertCollectionRequestWasteItem(context, gen.InsertCollectionRequestWasteItemParams{
			CollectionRequestID: params.CollectionRequestID,
			WasteTypeID:         int32(wasteItem.WasteTypeID),
			CollectorID:         int32(auth.UserCompanyId.Int64),
			Weight:              wasteItem.Weight,
		})
		if insertError != nil {
			errorInsertingWasteItem = true
		}
	}

	if errorInsertingWasteItem {
		gen.REPO.DeleteWasteItemsForCollectionRequest(context, params.CollectionRequestID)
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to insert waste items",
		})
		return
	}

	for _, v := range params.Waste {
		item, err := gen.REPO.GetInventoryItem(
			context, gen.GetInventoryItemParams{
				WasteTypeID: sql.NullInt32{Int32: int32(v.WasteTypeID), Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64)})
		if err != nil && err == sql.ErrNoRows {
			gen.REPO.InsertToInventory(context, gen.InsertToInventoryParams{
				TotalWeight: v.Weight,
				CompanyID:   int32(auth.UserCompanyId.Int64),
				WasteTypeID: sql.NullInt32{Int32: int32(v.WasteTypeID), Valid: true},
			})
		} else if err != nil && err != sql.ErrNoRows {
			//errorSavingInventory = true
			logger.Log("CollectionRequestWasteItemsController/InsertWasteItems", fmt.Sprint("Error saving to inventory :: ", err.Error()), logger.LOG_LEVEL_ERROR)
		} else {
			currentQuantity := item.TotalWeight
			var remainingWeight = currentQuantity + v.Weight
			//update with the remaining weight
			gen.REPO.UpdateInventoryItem(context, gen.UpdateInventoryItemParams{
				TotalWeight: remainingWeight,
				ID:          item.ID,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Inserted waste items",
	})

}

func (controller CollectionRequestsController) GetCollections(context *gin.Context) {

}
