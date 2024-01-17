package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)


type ChampionCollectorController struct{}

type AssignChampionToCollectorParams struct {
	ChampionID int32  `json:"champion_id"  binding:"required"`
	CollectorID int32  `json:"collector_id"  binding:"required"`
}

type UpdateChampionCollectorParams struct {
	CollectorID int32 `json:"collector_id" binding:"required"`
	ID          int32         `json:"id" binding:"required"`
}


func (championCollectorController ChampionCollectorController) AssignChampionToCollector(context *gin.Context) {
	var params AssignChampionToCollectorParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	count, err := gen.REPO.GetAssignedCollectorsToGreenChampion(
	context,sql.NullInt32{Int32: params.ChampionID, Valid: true});
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Champion Has an assigned Collector",
		})
		return
	}

	ChampionCollector, insertError := gen.REPO.AssignChampionToCollector(context, gen.AssignChampionToCollectorParams{
		ChampionID:		sql.NullInt32{Int32: params.ChampionID, Valid: true},
		CollectorID:    sql.NullInt32{Int32: params.CollectorID, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": insertError,
		})
		return
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Assigned Collector to Champion",
		"Assigned": ChampionCollector, // Include the company details in the response
	})

}


func (championCollectorController ChampionCollectorController) GetTheCollectorForAChampion(context *gin.Context) {
	// Retrieve champion ID from the URL parameter
	championIDParam := context.Param("id")

	// Convert the champion ID parameter to an int32
	championID, err := strconv.ParseInt(championIDParam, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid champion ID",
		})
		return
	}

	// Create a sql.NullInt32 instance for Championid
	championIDNullable := sql.NullInt32{Int32: int32(championID), Valid: true}

	// Call the repository method with the champion ID
	ChampionCollector, err := gen.REPO.GetTheCollectorForAChampion(context, championIDNullable)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":        false,
		"result": ChampionCollector,
	})
}

func (championCollectorController ChampionCollectorController) GetAllChampionsForACollector(context *gin.Context) {
	// Retrieve champion ID from the URL parameter
	collectorIDParam := context.Param("id")

	// Convert the champion ID parameter to an int32
	collectorID, err := strconv.ParseInt(collectorIDParam, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid collector ID",
		})
		return
	}

	// Create a sql.NullInt32 instance for Championid
	collectorIDNullable := sql.NullInt32{Int32: int32(collectorID), Valid: true}

	// Call the repository method with the champion ID
	Champions, err := gen.REPO.GetAllChampionsForACollector(context, collectorIDNullable)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":        false,
		"result": Champions,
	})
}


func (championCollectorController ChampionCollectorController) UpdateChampionCollector(context *gin.Context) {
	var params UpdateChampionCollectorParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}


	insertError := gen.REPO.UpdateChampionCollector(context, gen.UpdateChampionCollectorParams{
		ID:				params.ID,
		CollectorID:    sql.NullInt32{Int32: params.CollectorID, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": insertError,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated Champion Collector",
	})

}