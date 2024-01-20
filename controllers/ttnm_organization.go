package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type TtnmOrganizationController struct{}

type UpdateTtnmOrganizationProfileParams struct {
	OrganizationID         string `json:"organization_id"`
	Name                   string `json:"name"`
	TagLine                string `json:"tag_line"`
	AboutUs                string `json:"about_us"`
	LogoPath               string `json:"logo_path"`
	WebsiteUrl             string `json:"website_url"`
	City                   string `json:"city"`
	State                  string `json:"state"`
	Zip                    string `json:"zip"`
	Country                string `json:"country"`
	AppAppstoreLink        string `json:"app_appstore_link"`
	AppGooglePlaystoreLink string `json:"app_google_playstore_link"`
}

func (ttnmOrganizationController TtnmOrganizationController) GetTTNMOrganizations(context *gin.Context) {
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

func (ttnmOrganizationController TtnmOrganizationController) UpdateTtnmOrganizationProfile(context *gin.Context) {
	var params UpdateTtnmOrganizationProfileParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	updateError := gen.REPO.UpdateMainOrganizationProfile(context, gen.UpdateMainOrganizationProfileParams{
		OrganizationID: params.OrganizationID,
		Name:           params.Name,
		TagLine:        params.TagLine,
		AboutUs:        params.AboutUs,
		LogoPath:       params.LogoPath,
		WebsiteUrl:     params.WebsiteUrl,
		Country:        params.Country,
		City:           params.City,
		Zip:            params.Zip,
		State:          params.State,
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
		"message": "Successfully updated TTNM Profile",
	})
}
