package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "strings"
	"ttnmwastemanagementsystem/gen"
)

type PresetController struct{}

func (presetController PresetController) GetPresetValue(context *gin.Context) {
	query := context.Query("q")
	countries, _ := gen.REPO.GetAllCountries(context)
	organizations, _ := gen.REPO.GetAllOrganizations(context)

	if query == "countries" {
		context.JSON(http.StatusOK, gin.H{
			"error":     false,
			"countries": countries,
		})
	} else if query == "countries_and_organizations" {
		context.JSON(http.StatusOK, gin.H{
			"error":         false,
			"countries":     countries,
			"organizations": organizations,
		})
	}
}
