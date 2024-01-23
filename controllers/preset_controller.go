package controllers

import (
	"net/http"
	_ "strings"
	"ttnmwastemanagementsystem/gen"
	"github.com/gin-gonic/gin"
)

type PresetController struct{}

func (presetController PresetController) GetCountries(context *gin.Context) {
	countries,_:=gen.REPO.GetAllCountries(context)
	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"countries": countries,
	})
}
