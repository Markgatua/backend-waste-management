// controllers/waste_collection_controller.go

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

// GetWasteCollectionReports handles the GET request for waste collection reports
func GetWasteCollectionReports(c *gin.Context) {
	// Assuming you have a function to get a database connection named "GetDB"
	db := GetDB()
	defer db.Close()

	producerID := c.Param("producerID") // Assuming producerID is a parameter in the URL
	timeframe := c.Query("timeframe")  // Query parameter for timeframe (e.g., "monthly", "yearly", "weekly")
	wasteType := c.Query("wasteType")   // Query parameter for waste type

	// Construct the SQL query based on the provided parameters
	query := `
		SELECT 
			cr.id AS request_id,
			cr.producer_id,
			cr.collector_id,
			cr.request_date,
			cr.pickup_time_stamp_id,
			cr.location,
			cr.administrative_level_1_location,
			cr.lat,
			cr.lng,
			cr.pickup_date,
			cr.status,
			cr.first_contact_person,
			cr.second_contact_person,
			cr.created_at AS request_created_at,
			cri.id AS waste_item_id,
			cri.waste_type_id,
			cri.collector_id AS waste_collector_id,
			cri.weight,
			cri.created_at AS waste_item_created_at
		FROM 
			collection_requests cr
		LEFT JOIN 
			collection_request_waste_items cri ON cr.id = cri.collection_request_id
		LEFT JOIN 
			companies AS champion ON champion.id = cr.producer_id
		LEFT JOIN 
			companies AS collector ON collector.id = cr.collector_id
		WHERE 
			cr.producer_id = $1
			AND cr.status IN (1, 2, 3, 4, 5)
	`

	// Add additional conditions based on the provided parameters
	if timeframe != "" {
		query += fmt.Sprintf(" AND DATE_TRUNC('%s', cr.request_date) = DATE_TRUNC('%s', CURRENT_DATE)", timeframe, timeframe)
	}

	if wasteType != "" {
		query += fmt.Sprintf(" AND cri.waste_type_id = %s", wasteType)
	}

	rows, err := db.Query(query, producerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var wasteReports []models.WasteCollectionReport
	for rows.Next() {
		var wasteReport models.WasteCollectionReport
		err := rows.Scan(
			&wasteReport.RequestID,
			&wasteReport.ProducerID,
			&wasteReport.CollectorID,
			&wasteReport.RequestDate,
			&wasteReport.PickupTimeStampID,
			&wasteReport.Location,
			&wasteReport.AdminLevel1Location,
			&wasteReport.Lat,
			&wasteReport.Lng,
			&wasteReport.PickupDate,
			&wasteReport.Status,
			&wasteReport.FirstContactPerson,
			&wasteReport.SecondContactPerson,
			&wasteReport.RequestCreatedAt,
			&wasteReport.WasteItemID,
			&wasteReport.WasteTypeID,
			&wasteReport.WasteCollectorID,
			&wasteReport.Weight,
			&wasteReport.WasteItemCreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		wasteReports = append(wasteReports, wasteReport)
	}

	c.JSON(http.StatusOK, wasteReports)
}

// GetWasteCollectionReportChartAPI handles the GET request for waste collection report chart data
func GetWasteCollectionReportChartAPI(c *gin.Context) {
	// Add logic to fetch and format data for chart API based on your requirements
	// Example: Return JSON data for a chart
	chartData := map[string]interface{}{
		"labels":   []string{"January", "February", "March", "April", "May"},
		"datasets": []map[string]interface{}{{"label": "Monthly Report", "data": []int{10, 15, 20, 25, 30}}},
	}

	c.JSON(http.StatusOK, chartData)
}
