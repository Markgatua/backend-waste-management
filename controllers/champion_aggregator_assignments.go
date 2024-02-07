package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"
)

type ChampionCollectorController struct{}

type AssignChampionToCollectorParams struct {
	ChampionID  int32 `json:"champion_id"  binding:"required"`
	CollectorID int32 `json:"collector_id"  binding:"required"`
}

type AssignAggregatorsToGreenChampionsParams struct {
	Aggregators     []Aggregator `json:"aggregators" binding:"required"`
	GreenChampionID int32        `json:"green_champion_id" binding:"required"`
}

type Aggregator struct {
	ID          int64         `json:"id"`
	PickupTimes []PickUpTimes `json:"pickup_times"`
}

type PickUpTimes struct {
	PickupID  int32  `json:"pickup_id" binding:"required"`
	PickupDay string `json:"pickup_day"`
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
	if err != nil && err != sql.ErrNoRows {
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
	for _, v := range params.Aggregators {
		company, err := gen.REPO.GetCompany(context, int32(v.ID))
		if err != nil && err != sql.ErrNoRows {
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
	for _, v := range params.Aggregators {
		value, err := gen.REPO.AssignCollectorsToGreenChampion(context, gen.AssignCollectorsToGreenChampionParams{
			ChampionID:  params.GreenChampionID,
			CollectorID: int32(v.ID),
			//PickupDay:   null.StringFrom(v.PickupDay).NullString,
			//PickupTime:  null.StringFrom(v.PickupTime).NullString,
		})
		if err != nil {
			logger.Log("ChamptionAggregatorAssignmentController", err.Error(), logger.LOG_LEVEL_ERROR)
		}
		for _, x := range v.PickupTimes {
			err = gen.REPO.SetPickupTimesForGreenChampion(context, gen.SetPickupTimesForGreenChampionParams{
				ChampionAggregatorAssignmentID: value.ID,
				PickupTimeStampID:              x.PickupID,
				PickupDay:                      x.PickupDay,
			})
			if err != nil {
				logger.Log("ChamptionAggregatorAssignmentController [Pickup times]", err.Error(), logger.LOG_LEVEL_ERROR)
			}
		}
	}
	context.JSON(http.StatusOK, gin.H{
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

func (championCollectorController ChampionCollectorController) GetCollectorsForGreenChampion(context *gin.Context) {
	// Retrieve champion ID from the URL parameter
	//
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	greenChampionID := context.Query("gid")

	searchQuery := ""
	companyQuery := ""

	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.collector_company_name ilike " + "'%" + search + "%'" + " or q.collector_company_location ilike " + "'%" + search + "%'" + " or q.collector_contact_person1_first_name ilike " + "'%" + search + "%'" + " or q.collector_contact_person1_phone ilike " + "'%" + search + "%'" + ")"
	}

	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if greenChampionID == "" {
		greenChampionID = fmt.Sprint(auth.UserCompanyId.Int64)
	}

	companyQuery = " and q.champion_id=" + greenChampionID

	// Retrieve champion ID from the URL parameter
	greenChampionIDParam := context.Param("id")

	query := `
	 select * from
	 (
		select
		champion_aggregator_assignments.id,
		champion_aggregator_assignments.champion_id,
		champion_aggregator_assignments.collector_id,
		companies.name as  champion_company_name,
		companies.location as champion_company_location,
		champion_aggregator_assignments.created_at,

		companies.contact_person1_first_name as collector_contact_person1_first_name,
        companies.contact_person1_last_name as collector_contact_person1_last_name,
        companies.contact_person1_phone as collector_contact_person1_phone,
        companies.contact_person1_email as collector_contact_person1_email,
        companies.contact_person2_email as collector_contact_person2_email,
        companies.contact_person2_first_name as collector_contact_person2_first_name,
        companies.contact_person2_last_name as collector_contact_person2_last_name,
        companies.contact_person2_phone as collector_contact_person2_phone

		from champion_aggregator_assignments

		inner join companies on companies.id=champion_aggregator_assignments.collector_id

	 ) as q where q.created_at is not null and q.champion_id=` + greenChampionIDParam + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from champion_aggregator_assignments where created_at is not null and champion_id=", greenChampionIDParam))

	if err != nil {
		logger.Log("AggregatorController/GetCollectorsForGreenChampion", err.Error(), logger.LOG_LEVEL_ERROR)
	}
	logger.Log("AggregatorController/GetCollectorsForGreenChampion", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range results {
		assignmentID, _ := v["id"]
		query := fmt.Sprint(`
		select champion_pickup_times.*,
		pickup_time_stamps.stamp,
		pickup_time_stamps.time_range
		from champion_pickup_times left join pickup_time_stamps on pickup_time_stamps.id=champion_pickup_times.pickup_time_stamp_id

		where champion_pickup_times.champion_aggregator_assignment_id=
		`, assignmentID)

		pickupTimes, err := utils.Select(gen.REPO.DB, query)

		if err != nil {
			logger.Log("AggregatorController/GetCollectorsForGreenChampion[pickup time]", query, logger.LOG_LEVEL_ERROR)
		}

		//if pickupTimes==nil{
		//	pickupTimes= []any{}
		//}

		v["pickup_times"] = pickupTimes
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"total_count": totalCount,
		"content":     results,
	})
}

func (championCollectorController ChampionCollectorController) GetAllChampionsForACollector(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")

	searchQuery := ""
	companyQuery := ""

	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.champion_company_name ilike " + "'%" + search + "%'" + " or q.champion_company_location ilike " + "'%" + search + "%'" + " or q.champion_contact_person1_first_name ilike " + "'%" + search + "%'" + " or q.champion_contact_person1_phone ilike " + "'%" + search + "%'" + ")"
	}

	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if companyID == "" {
		companyID = fmt.Sprint(auth.UserCompanyId.Int64)
	}

	companyQuery = " and q.collector_id=" + companyID

	// Retrieve champion ID from the URL parameter
	collectorIDParam := context.Param("id")

	query := `
	 select * from
	 (
		select
		champion_aggregator_assignments.id,
		champion_aggregator_assignments.champion_id,
		champion_aggregator_assignments.collector_id,
		companies.name as  champion_company_name,
		companies.location as champion_company_location,
		champion_aggregator_assignments.created_at,

		companies.contact_person1_first_name as champion_contact_person1_first_name,
        companies.contact_person1_last_name as champion_contact_person1_last_name,
        companies.contact_person1_phone as champion_contact_person1_phone,
        companies.contact_person1_email as champion_contact_person1_email,
        companies.contact_person2_email as champion_contact_person2_email,
        companies.contact_person2_first_name as champion_contact_person2_first_name,
        companies.contact_person2_last_name as champion_contact_person2_last_name,
        companies.contact_person2_phone as champion_contact_person2_phone

		from champion_aggregator_assignments

		inner join companies on companies.id=champion_aggregator_assignments.collector_id

	 ) as q where q.created_at is not null and q.collector_id=` + collectorIDParam + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from champion_aggregator_assignments where created_at is not null and collector_id=", collectorIDParam))

	if err != nil {
		logger.Log("AggregatorController/GetAllChampionsForACollector", err.Error(), logger.LOG_LEVEL_ERROR)
	}
	logger.Log("AggregatorController/GetAllChampionsForACollector", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range results {
		assignmentID, _ := v["id"]
		query := fmt.Sprint(`
		select champion_pickup_times.*,
		pickup_time_stamps.stamp,
		pickup_time_stamps.time_range
		from champion_pickup_times left join pickup_time_stamps on pickup_time_stamps.id=champion_pickup_times.pickup_time_stamp_id

		where champion_pickup_times.champion_aggregator_assignment_id=
		`, assignmentID)

		pickupTimes, err := utils.Select(gen.REPO.DB, query)

		if err != nil {
			logger.Log("AggregatorController/GetAllChampionsForACollector[pickup time]", query, logger.LOG_LEVEL_ERROR)
		}

		//if pickupTimes==nil{
		//	pickupTimes= []any{}
		//}

		v["pickup_times"] = pickupTimes
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"total_count": totalCount,
		"content":     results,
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
