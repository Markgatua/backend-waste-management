package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type ChampionCollectorController struct{}

type AssignChampionToCollectorParams struct {
	ChampionID  int32 `json:"champion_id"  binding:"required"`
	CollectorID int32 `json:"collector_id"  binding:"required"`
}

type AssignAggregatorsToGreenChampionsParams struct {
	AggregatorIDs   []int32 `json:"aggregator_ids" binding:"required"`
	GreenChampionID int32   `json:"green_champion_id" binding:"required"`
}

type UpdateChampionCollectorParams struct {
	CollectorID int32 `json:"collector_id" binding:"required"`
	ID          int32 `json:"id" binding:"required"`
}

func (championCollectorController ChampionCollectorController) AssignAggregatorsToGreenChampionsParam(context *gin.Context) {
	var params AssignAggregatorsToGreenChampionsParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	greenChampion, err := gen.REPO.GetCompany(context, params.GreenChampionID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if greenChampion.CompanyType != 1 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "The selected green champion is not a green champion",
		})
		return
	}

	var error_ string = ""

	for _, v := range params.AggregatorIDs {
		company, err := gen.REPO.GetCompany(context, v)
		if err != nil {
			error_ = "Error getting aggregator"
		} else {
			if company.CompanyType != 2 {
				error_ = "One of the aggregators is not an aggregator"
			}
		}
	}
	if error_ != "" {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": error_,
		})
		return
	}

	err = gen.REPO.RemoveAggrigatorsAssignedFromGreenChampions(context, params.GreenChampionID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range params.AggregatorIDs {
		gen.REPO.AssignCollectorsToGreenChampion(context, gen.AssignCollectorsToGreenChampionParams{
			ChampionID:  params.GreenChampionID,
			CollectorID: v,
		})
	}

	context.JSON(http.StatusUnprocessableEntity, gin.H{
		"error":   false,
		"message": "Successfully assigned collectors to aggregator",
	})
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
		context, params.ChampionID)
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Champion Has an assigned Collector",
		})
		return
	}

	ChampionCollector, insertError := gen.REPO.AssignChampionToCollector(context, gen.AssignChampionToCollectorParams{
		ChampionID:  params.ChampionID,
		CollectorID: params.CollectorID,
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
		"error":    false,
		"message":  "Successfully Assigned Collector to Champion",
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

	// Call the repository method with the champion ID
	ChampionCollector, err := gen.REPO.GetTheCollectorForAChampion(context, int32(championID))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":  false,
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

	// Call the repository method with the champion ID
	Champions, err := gen.REPO.GetAllChampionsForACollector(context, int32(collectorID))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":  false,
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
		ID:          params.ID,
		CollectorID: params.CollectorID,
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
