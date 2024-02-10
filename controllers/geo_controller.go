package controllers

import (
	"net/http"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/utils"

	"github.com/gin-gonic/gin"
)

type GeoController struct{}

func (c GeoController) GetAllCountries(context *gin.Context) {
	countries, _ := utils.Select(gen.REPO.DB,"select * from countries");//  gen.REPO.GetAllCountries(context)


	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": countries, // Include the company details in the response
	})
}
