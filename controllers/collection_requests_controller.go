package controllers

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
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
			logger.Log("CollectionRequestWasteItemsController", insertError.Error(), logger.LOG_LEVEL_ERROR)
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

func (aggregatorController CollectionRequestsController) GetCollectionSchedule(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	dateRangeStart := context.Query("sd")
	search := context.Query("s")
	dateRangeEnd := context.Query("ed")
	companyID := context.Query("cid")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	status := context.Query("st") //1- Pending, 2- confirmed, 3- on the way, 4 cancelled, 5 completed

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	limitOffset := ""
	statusQuery := ""

	if search != "" {
		searchQuery = " and (q.time_range ilike " + "'%" + search + "%'" + " or q.champion_name ilike " + "'%" + search + "%'" + ")"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if companyID == "" {
		companyID = fmt.Sprint(auth.UserCompanyId.Int64)
		//companyQuery = fmt.Sprint(" and  q.company_id=", auth.UserCompanyId.Int64)
	}
	if status != "" {
		statusQuery = " and q.status=" + status
	}
	companyQuery = " and q.collector_id=" + companyID

	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.request_date as date)>='" + dateRangeStart + "' and cast(q.request_date as date)<='" + dateRangeEnd + "'"
	}

	query := `
	 select * from 
	 (
		select 
		collection_requests.id,
		collection_requests.producer_id,
		companies.name as champion_name,
		collection_requests.collector_id,
		collection_requests.request_date,
		collection_requests.pickup_date,
		collection_requests.status,
		collection_requests.location as pickup_location,

		collection_requests.lat,
		collection_requests.lng,

		collection_requests.created_at,
		collection_requests.pickup_time_stamp_id,
		collection_requests.id,
		collection_requests.first_contact_person,
		collection_requests.second_contact_person,
		pickup_time_stamps.stamp,
		pickup_time_stamps.position,
		pickup_time_stamps.time_range

		from collection_requests 

		inner join pickup_time_stamps on pickup_time_stamps.id=collection_requests.pickup_time_stamp_id
		inner join companies on companies.id=collection_requests.producer_id
	 ) as q where q.created_at is not null` + dateRangeQuery + statusQuery + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	// var totalCount = 0
	// err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from collection_request_waste_items where created_at is not null and collector_id=", companyID))

	//fmt.Println(err.Error())
	logger.Log("CollectionRequestsController/GetCollections", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var addRequestsToSchedule bool = true
	now := carbon.Now()
	group := make(map[string][]interface{})
	for _, v := range results {
		v["type"] = "Collection request"
		requestDate, _ := v["request_date"]

		if addRequestsToSchedule {
			dayOfWeek := ""
			carbonTime := carbon.Parse(requestDate.(time.Time).Format("2006-01-02"))
			if carbonTime.IsMonday() {
				dayOfWeek = "Monday"
			} else if carbonTime.IsTuesday() {
				dayOfWeek = "Tuesday"
			} else if carbonTime.IsWednesday() {
				dayOfWeek = "Wednesday"
			} else if carbonTime.IsThursday() {
				dayOfWeek = "Thursday"
			} else if carbonTime.IsFriday() {
				dayOfWeek = "Friday"
			} else if carbonTime.IsSaturday() {
				dayOfWeek = "Saturday"
			} else if carbonTime.IsSunday() {
				dayOfWeek = "Sunday"
			}

			if carbonTime.DiffInWeeks(now) == 0 {
				items, ok := group[dayOfWeek]
				if ok {
					items = append(items, v)
					group[dayOfWeek] = items
				} else {
					var values []interface{}
					values = append(values, v)
					group[dayOfWeek] = values
				}
			}

			//fmt.Println(dayOfWeek)
		} else {
			//date, _ := time.Parse("2006-01-02", requestDate.(string))
			key := requestDate.(time.Time).Format("2006-01-02")
			items, ok := group[key]
			if ok {
				items = append(items, v)
				group[key] = items
			} else {
				var values []interface{}
				values = append(values, v)
				group[key] = values
			}
		}
	}

	collectionScheduleQuery := `
	select * from 
	(
	   select 
	   champion_pickup_times.id as pickup_time_id,
	   champion_pickup_times.champion_aggregator_assignment_id,
	   champion_pickup_times.pickup_time_stamp_id,
	   champion_pickup_times.pickup_day,
	   champion_aggregator_assignments.champion_id,
	   champion_aggregator_assignments.collector_id,

	   companies.lat,
	   companies.lng,
	   companies.name as champion_name,
	   companies.location,

	   companies.contact_person1_first_name,
       companies.contact_person1_last_name,
       companies.contact_person1_phone,
       companies.contact_person1_email,
       companies.contact_person2_email,
	   companies.administrative_level_1_location,
       companies.contact_person2_first_name,
       companies.contact_person2_last_name,
       companies.contact_person2_phone,
	   companies.location as pickup_location,


	   pickup_time_stamps.stamp,
	   pickup_time_stamps.time_range,
	   pickup_time_stamps.position


	   from champion_pickup_times 

	   left join pickup_time_stamps on pickup_time_stamps.id=champion_pickup_times.pickup_time_stamp_id
	   left join champion_aggregator_assignments on champion_pickup_times.champion_aggregator_assignment_id=champion_aggregator_assignments.id
	   left join companies on companies.id = champion_aggregator_assignments.champion_id

	) as q where q.pickup_time_id is not null` //+ dateRangeQuery + statusQuery + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	collectionScheduleResults, err := utils.Select(gen.REPO.DB, collectionScheduleQuery)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range collectionScheduleResults {
		key, _ := v["pickup_day"]
		items, ok := group[key.(string)]
		if ok {
			items = append(items, v)
			group[key.(string)] = items
		} else {
			var values []interface{}
			values = append(values, v)
			group[key.(string)] = values
		}
		v["type"] = "Collection schedule"
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": group,
	})
}

func (aggregatorController CollectionRequestsController) ChangeCollectionRequestStatus(context *gin.Context) {
	collectionRequestID := context.Param("id")
	status := context.Param("status")

	var id32 int32
	fmt.Sscan(collectionRequestID, &id32)

	var status_ int32
	fmt.Sscan(status, &status_)

	if !utils.InArray(status, []string{"1", "2", "3", "4", "5"}) {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid status",
		})
		return
	}
	_, err := gen.REPO.GetCollectionRequest(context, id32)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = gen.REPO.ChangeCollectionRequestStatus(context, gen.ChangeCollectionRequestStatusParams{
		Status: status_,
		ID:     id32,
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Successfully updated collection request status",
	})

}

func (aggregatorController CollectionRequestsController) GetAggregatorCollectionRequests(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	status := context.Query("st")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")
	dateRangeStart := context.Query("sd")
	dateRangeEnd := context.Query("ed")

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	statusQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.champion_company_name ilike " + "'%" + search + "%'" + " or q.location ilike " + "'%" + search + "%'" + " or q.stamp ilike " + "'%" + search + "%'" + ")"
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

	if status != "" {
		statusQuery = " and q.status=" + status
	}

	companyQuery = " and q.collector_id=" + companyID

	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.request_date as date)>='" + dateRangeStart + "' and cast(q.request_date as date)<='" + dateRangeEnd + "'"
	}
	query := `
	 select * from 
	 (
		select 
		collection_requests.id,
		collection_requests.request_date,
		collection_requests.pickup_date,
		collection_requests.collector_id,
		collection_requests.producer_id,
		collection_requests.status,
		collection_requests.producer_id,
		collection_requests.pickup_time_stamp_id,
		collection_requests.location,
		collection_requests.administrative_level_1_location,
		collection_requests.lat,
		collection_requests.lng,
		collection_requests.first_contact_person,
		collection_requests.second_contact_person,
		collection_requests.created_at,
        pickup_time_stamps.stamp,
	    pickup_time_stamps.time_range,
	    pickup_time_stamps.position,
		companies.name as champion_company_name

		from collection_requests 

		inner join companies on companies.id=collection_requests.producer_id
		left join pickup_time_stamps on pickup_time_stamps.id=collection_requests.pickup_time_stamp_id

	 ) as q where q.created_at is not null` + dateRangeQuery + searchQuery + companyQuery + statusQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from collection_requests where created_at is not null and collector_id=", companyID))

	//fmt.Println(err.Error())
	logger.Log("CollectionRequestsController/GetAggregatorCollectionRequests", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range results {
		key, _ := v["id"]
		query = fmt.Sprint(`select collection_request_waste_items.*,waste_types.name from collection_request_waste_items inner join waste_types on waste_types.id=collection_request_waste_items.waste_type_id where collection_request_id=`, key)
		//collection_request_waste_items

		//fmt.Print(query)
		results, err := utils.Select(gen.REPO.DB, query)
		if err != nil {
			logger.Log("CollectionRequestsController/GetAggregatorCollectionRequests/[wasteitems]", err.Error(), logger.LOG_LEVEL_ERROR)
		}
		v["collection_waste_items"] = results
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}
