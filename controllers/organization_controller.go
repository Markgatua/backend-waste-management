package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

type OrgnizationController struct{}

type InsertOrganizationParam struct {
	Name             string `json:"name"  binding:"required"`
	CountryID        int32  `json:"country_id"  binding:"required"`
	OrganizationType int32  `json:"organization_type" binding:"required"`
	LogoPath         string `json:"logo_path"`

	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
type SetActiveInactiveOrganizationStatusParam struct {
	OrganizationID int64 `json:"organization_id" binding:"required"`
	IsActive       *bool `json:"is_active" binding:"required"`
}

type UpdateOrganizationParams struct {
	ID               int    `json:"id"  binding:"required"`
	Name             string `json:"name"  binding:"required"`
	CountryID        int32  `json:"country_id"  binding:"required"`
	OrganizationType int32  `json:"organization_type" binding:"required"`
	LogoPath         string `json:"logo_path"`

	UserID    int32  `json:"user_id" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password"`
}

func (controller OrgnizationController) SetActiveInActiveStatus(context *gin.Context) {
	var param SetActiveInactiveOrganizationStatusParam
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err = gen.REPO.UpdateOrganizationIsActive(context, gen.UpdateOrganizationIsActiveParams{
		IsActive: *param.IsActive,
		ID:       int32(param.OrganizationID),
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error setting user status",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Updated user status",
		})
	}
}

func (c OrgnizationController) GetPresets(context *gin.Context) {
	countries, _ := gen.REPO.GetAllCountries(context)
	context.JSON(http.StatusOK, gin.H{
		"error":     false,
		"countries": countries,
	})

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

	user, err := GetEmailUser(params.Email)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error adding organization",
		})
		return
	}
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email already exists`",
		})
		return
	}
	ogranization, insertError := gen.REPO.InsertOrganization(context, gen.InsertOrganizationParams{
		Name:             params.Name,
		CountryID:        params.CountryID,
		OrganizationType: params.OrganizationType,
	})

	var roleID = -1
	var userType = -1
	if params.OrganizationType == 1 {
		roleID = 2
		userType = 2
	} else if params.OrganizationType == 2 {
		roleID = 6
		userType = 8
	}

	_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,provider,role_id,user_organization_id,user_type,created_at,updated_at,password,is_organization_super_admin) VALUES (:email,:first_name,:last_name,:provider,:role_id,:user_organization_id,:user_type,:created_at,:updated_at,:password,:is_organization_super_admin)`,
		map[string]interface{}{
			"email":                       params.Email,
			"first_name":                  params.FirstName,
			"last_name":                   params.LastName,
			"provider":                    "email",
			"role_id":                     roleID,
			"user_organization_id":        ogranization.ID,
			"is_organization_super_admin": true,
			"user_type":                   userType,
			"password":                    helpers.Functions{}.HashPassword(params.Password),
			"created_at":                  time.Now(),
			"updated_at":                  time.Now(),
		})

	if params.LogoPath != "" {
		UploadController{}.SaveToUploadsTable(params.LogoPath, "organizations", ogranization.ID)
	}
	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add organization",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully created organization",
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

	usersWithEmailWithoutID, err := gen.REPO.GetUserWithEmailWithoutID(context, gen.GetUserWithEmailWithoutIDParams{
		Email: sql.NullString{String: params.Email, Valid: true},
		ID:    params.UserID,
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if len(usersWithEmailWithoutID) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another user has the same email",
		})
		return
	}

	// Update Waste Group
	updateError := gen.REPO.UpdateOrganization(context, gen.UpdateOrganizationParams{
		Name:             params.Name,
		CountryID:        int32(params.CountryID),
		OrganizationType: params.OrganizationType,
		ID:               int32(params.ID),
	})
	var roleID = -1
	var userType = -1
	if params.OrganizationType == 1 {
		roleID = 2
		userType = 2
	} else if params.OrganizationType == 2 {
		roleID = 6
		userType = 8
	}

	if params.Password == "" {
		gen.REPO.UpdateUserFirstNameLastNameEmailRoleAndUserType(context, gen.UpdateUserFirstNameLastNameEmailRoleAndUserTypeParams{
			Email:     null.StringFrom(params.Email).NullString,
			RoleID:    sql.NullInt32{Int32: int32(roleID), Valid: true},
			UserType:  sql.NullInt16{Int16: int16(userType), Valid: true},
			FirstName: sql.NullString{String: params.FirstName, Valid: true},
			LastName:  sql.NullString{String: params.LastName, Valid: true},
			ID:        params.UserID,
		})

	} else {
		gen.REPO.UpdateUserFirstNameLastNameEmailRoleUserTypeAndPassword(context, gen.UpdateUserFirstNameLastNameEmailRoleUserTypeAndPasswordParams{
			Email:     null.StringFrom(params.Email).NullString,
			RoleID:    sql.NullInt32{Int32: int32(roleID), Valid: true},
			UserType:  sql.NullInt16{Int16: int16(userType), Valid: true},
			ID:        params.UserID,
			FirstName: sql.NullString{String: params.FirstName, Valid: true},
			LastName:  sql.NullString{String: params.LastName, Valid: true},
			Password:  null.StringFrom(helpers.Functions{}.HashPassword(params.Password)).NullString,
		})

	}

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
