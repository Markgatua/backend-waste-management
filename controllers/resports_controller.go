package controllers

import (
	"context"
	"net/http"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type ReportsController struct{}

func (controller ReportsController) GetMainOrganizationStats(context *gin.Context) {
	var startDate = context.Query("start_date")
	var endDate = context.Query("end_date")

	organizationStats, err := GetOrganizationCount(startDate, endDate)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	branchStats, err := GetBranchCount(startDate, endDate)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"content": gin.H{
			"organization": organizationStats,
			"branches":     branchStats,
		},
	})
}

func GetOrganizationCount(startDate string, endDate string) ([]gen.GetOrganizationCountRow, error) {
	return gen.REPO.GetOrganizationCount(context.Background())
}

func GetBranchCount(startDate string, endDate string) ([]gen.GetBranchCountRow, error) {
	return gen.REPO.GetBranchCount(context.Background())
}

func GetMainSystemUsersCount(startDate string, endDate string) {

}

func GetOrganizationUsersCount(startDate string, endDate string) {

}

func GetCollectionStats(startDate string, endDate string) {

}

func GetWasteTypeStats(startDate string, endDate string) {

}

func GetCollectionStartsByLocation(startDate string, endDate string) {

}

// sales reports
func SalesOverTimeByBuyer() {

}
func SalesOverTimeByItem() {

}
func SalesOverTimeMoneyIn() {

}
func StockMovementStatement() {

}

//end of sales report

// Purchases reports
func PurchasesOverTimeBySupplier() {

}
func PurchasesOverTimeByItem() {

}
func PurchasesMadeOverTimeMoneyOut() {

}

//end of purchases report

// Collection reports
func CollectionOverTimeByCustomer() {

}
func CollectionOverTimeByLocation() {

}
func CollectionOverTimeByItem() {

}

//end of collection reports

// Dashboard
func SalesMadeOverTime() {
	//line graph
}

func StockLevelsByItem() {
	//bar chart
}

func CollectionStatus() {
	// completed, cancelled,
}

func InventoryStatus() {
	//low stock, in stock
}
//Dashboard