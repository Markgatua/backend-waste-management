package controllers

import (
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type OrgnizationController struct{}

type InsertOrganizationParam struct {
	Name      string `json:"name"  binding:"required"`
	CountryID int32  `json:"country_id"  binding:"required"`
}

type UpdateOrganizationParams struct {
	ID        int    `json:"id"  binding:"required"`
	CountryID int    `json:"country_id" binding:"required"`
	Name      string `json:"name"  binding:"required"`
}

func (c OrgnizationController) InsertOrganization(context *gin.Context) {
	var params InsertOrganizationParam
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	count, err := gen.REPO.GetOrganizationCountWithNameAndCountry(context, gen.GetOrganizationCountWithNameAndCountryParams{
		Name:      strings.ToLower(params.Name),
		CountryID: params.CountryID,
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Organization with the same name in the country already exists",
		})
		return
	}

	ogranization, insertError := gen.REPO.InsertOrganization(context, gen.InsertOrganizationParams{
		Name:      params.Name,
		CountryID: params.CountryID,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Company Region",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Company Region",
		"content": ogranization,
	})
}

func (c OrgnizationController) GetAllOrganizations(context *gin.Context) {
	organizations, err := gen.REPO.GetAllOrganizations(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": organizations,
	})
}

func (c OrgnizationController) GetOrganization(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	organization, err := gen.REPO.GetOrganization(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":        false,
		"organization": organization,
	})
}

func (c OrgnizationController) DeleteOrganization(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	err := gen.REPO.DeleteOrganization(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Organization successfully deleted",
	})
}

func (c OrgnizationController) UpdateOrganization(context *gin.Context) {
	var params UpdateOrganizationParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	count, err := gen.REPO.GetDuplicateOrganization(context, gen.GetDuplicateOrganizationParams{
		Name:      strings.ToLower(params.Name),
		CountryID: int32(params.CountryID),
		ID:        int32(params.ID),
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Organization with the same name in the country already exists",
		})
		return
	}

	// Update Waste Group
	updateError := gen.REPO.UpdateOrganization(context, gen.UpdateOrganizationParams{
		Name:      params.Name,
		CountryID: int32(params.CountryID),
		ID:        int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated Organization",
	})
}
