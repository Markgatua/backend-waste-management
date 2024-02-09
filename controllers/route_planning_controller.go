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
		CompanyID  int32      `json:"company_id" binding:"required"`
		Duration   float64    `json:"duration" binding:"required"`
		VehicleIDs []int32    `json:"vehicle_ids" binding:"required"`
		Shipments  []Shipment `json:"shipments" binding:"required"`
	}

	type Agents struct {
		StartLocation  []float64 `json:"start_location"`  //lat,lng
		EndLocation    []float64 `json:"end_location"`    //lat,lng
		PickupCapacity float64   `json:"pickup_capacity"` //1000 liters
	}
	type Job struct {
		Location     []float64 `json:"location"`
		Duration     int64     `json:"duration"`
		PickupAmount int64     `json:"pickup_amount"`
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
	company, err := gen.REPO.GetCompany(context, params.CompanyID)

	var agents []Agents
	var jobs []Job

	type BodyContent struct {
		Mode   string   `json:"mode"`
		Agents []Agents `json:"agents"`
		Jobs   []Job    `json:"jobs"`
	}

	bodyContent := BodyContent{}
	bodyContent.Mode = "drive"

	for _, v := range params.VehicleIDs {
		vehicle, _ := gen.REPO.GetVehicle(context, v)
		agents = append(agents, Agents{
			StartLocation:  []float64{company.Lat.Float64, company.Lng.Float64},
			EndLocation:    []float64{company.Lat.Float64, company.Lng.Float64},
			PickupCapacity: vehicle.Liters.Float64,
		})
	}

	for _, collectionSchedule := range collectionSchedules {
		jobs = append(jobs, Job{
			Location:     []float64{collectionSchedule.Lat.Float64, collectionSchedule.Lng.Float64},
			Duration:     int64(params.Duration),
			PickupAmount: 100,
		})
	}

	for _, collectionRequest := range collectionRequests {
		jobs = append(jobs, Job{
			Location:     []float64{collectionRequest.Lat.Float64, collectionRequest.Lng.Float64},
			Duration:     int64(params.Duration),
			PickupAmount: 100,
		})
	}

	bodyContent.Agents = agents
	bodyContent.Jobs = jobs

	context.JSON(http.StatusOK, gin.H{
		"error":                false,
		"collection_requests":  collectionRequests,
		"collection_schedules": collectionSchedules,
	})

}
