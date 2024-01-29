package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ttnmwastemanagementsystem/gen"
)

type GeoController struct{}

func (c GeoController) GetAllCountries(context *gin.Context) {
	countries, _ := gen.REPO.GetAllCountries(context)
	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": countries, // Include the company details in the response
	})
}
