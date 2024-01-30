package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type StatsController struct{}

func (controller StatsController) GetTTNMOrganizations(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	ttnm, err := gen.REPO.GetMainOrganization(context, id)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"Profile": ttnm,
	})
}

func GetOrganizationCount(startDate string, endDate string) {

}

func GetBranchCount(startDate string, endDate string) {

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
