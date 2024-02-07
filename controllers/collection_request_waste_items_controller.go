package controllers

import (
	"database/sql"
	_ "fmt"
	"net/http"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"

	"github.com/gin-gonic/gin"
)

type CollectionRequestWasteItemsController struct{}

type InsertWasteItemParams struct {
	CollectionRequestID int32   `json:"collection_request_id"`
	Waste               []Waste `json:"waste"`
}
type Waste struct {
	WasteTypeID int64   `json:"waste_type_id"`
	Weight      float64 `json:"weight"`
}

func (controller CollectionRequestWasteItemsController) InsertWasteItem(context *gin.Context) {
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
			Weight:              wasteItem.Weight,
		})
		if insertError != nil {
			errorInsertingWasteItem = true
		}
		// if insertError != nil {
		// 	context.JSON(http.StatusUnprocessableEntity, gin.H{
		// 		"error":           true,
		// 		"message":         "Failed to insert waste items",
		// 		"Collected Waste": collectedWaste,
		// 	})
		// 	return
		// }
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
		item, get_error := gen.REPO.GetInventoryItem(
			context, gen.GetInventoryItemParams{
				WasteTypeID: sql.NullInt32{Int32: int32(v.WasteTypeID), Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64)})

		if get_error != nil {
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
