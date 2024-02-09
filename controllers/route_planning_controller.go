package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "strings"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type RoutePlanningController struct{}

func (controller RoutePlanningController) GetRoutes(context *gin.Context) {
	logger.Log("RoutePlanningController", fmt.Sprint("Running route planner"), logger.LOG_LEVEL_INFO)

	type Shipment struct {
		ID                  int32 `json:"id" binding:"required"`
		IsCollectionRequest *bool `json:"is_collection_request" binding:"required"`
	}
	type Params struct {
		VehicleIDs []int32    `json:"vehicle_ids" binding:"required"`
		Shipments  []Shipment `json:"shipments" binding:"required"`
	}

	var params Params
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var collectionRequestsInArray []int32
	var collectionScheduleInArray []int32

	for _, v := range params.Shipments {
		if *v.IsCollectionRequest {
			//get lats and longs
			collectionRequestsInArray = append(collectionRequestsInArray, v.ID)
		} else {
			collectionScheduleInArray = append(collectionScheduleInArray, v.ID)
			//get lats and longs
		}
	}

	collectionRequests, err := gen.REPO.GetCollectionRequestsInArray(context, collectionRequestsInArray)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	collectionSchedules, err := gen.REPO.GetCollectionScheduleInArray(context, collectionScheduleInArray)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":               false,
		"collection_requests": collectionRequests,
		"collection_schedules": collectionSchedules,
	})

}
