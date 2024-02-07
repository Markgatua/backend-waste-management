package controllers

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"

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

func (aggregatorController CollectionRequestsController) GetCollections(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")
	dateRangeStart := context.Query("sd")
	dateRangeEnd := context.Query("ed")

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.waste_name ilike " + "'%" + search + "%'" + " or q.company_name ilike " + "'%" + search + "%'" + " or q.last_name ilike " + "'%" + search + "%'" + " or q.ref ilike " + "'%" + search + "%'" + ")"
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
	//else {
	//	companyQuery = " and  q.company_id=" + companyID
	//}

	companyQuery = " and q.company_id=" + companyID

	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.created_at as date)>='" + dateRangeStart + "' and cast(q.created_at as date)<='" + dateRangeEnd + "'"
	}
	query := `
	 select * from 
	 (
		select 
		collection_request_waste_items.id,
		collection_request_waste_items.waste_type_id,
		collection_request_waste_items.weight,
		collection_request_waste_items.collector_id,
		collection_request_waste_items.created_at,
		waste_types.name as waste_name,
		companies.name as company_name,
		from collection_request_waste_items 

		inner join companies on companies.id=collection_request_waste_items.collector_id
		inner join collection_request_waste_items on collection_request_waste_items.waste_type_id=waste_types.id
	 ) as q where q.created_at is not null` + dateRangeQuery + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from collection_request_waste_items where created_at is not null and collector_id=",companyID))

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

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}
