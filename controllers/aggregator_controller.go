package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/models"
	"ttnmwastemanagementsystem/utils"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type AggregatorController struct{}

type CreateAggregatorParams struct {
	Name           string          `json:"name"  binding:"required"`
	Location       models.Location `json:"location" binding:"required"`
	Email          string          `json:"email" binding:"required"`
	FirstName      string          `json:"first_name" binding:"required"`
	LastName       string          `json:"last_name" binding:"required"`
	Password       string          `json:"password" binding:"required"`
	OrganizationID int32           `json:"organization_id"`
	LogoPath       string          `json:"logo_path"`
	IsActive       *bool           `json:"is_active"  binding:"required"`
	Region         string          `json:"region"`
}

type UpdateAggregatorDataParams struct {
	Location       models.Location `json:"location" binding:"required"`
	Name           string          `json:"name"  binding:"required"`
	OrganizationID int32           `json:"organization_id"`
	IsActive       *bool           `json:"is_active"  binding:"required"`
	ID             int64           `json:"id"  binding:"required"`
	LogoPath       string          `json:"logo_path"`
	Region         string          `json:"region"`
	UserID         int32           `json:"user_id" binding:"required"`
	Email          string          `json:"email" binding:"required"`
	FirstName      string          `json:"first_name" binding:"required"`
	LastName       string          `json:"last_name" binding:"required"`
	Password       string          `json:"password"`
}

type UpdateAggregatorStatusParams struct {
	ID       int   `json:"id"  binding:"required"`
	IsActive *bool `json:"status"  binding:"required"`
}

func (controller AggregatorController) SetWasteTypes(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type WasteType struct {
		ID         int32   `json:"id" binding:"id"`
		AlertLevel float64 `json:"alert_level"`
	}
	type Params struct {
		//AggregatorID int32   `json:"aggregator_id"`
		WasteTypes []WasteType `json:"waste_types" binding:"required"`
	}
	var AggregatorID int32 = int32(auth.UserCompanyId.Int64)
	var params Params
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = gen.REPO.DeleteAggregatorWasteTypes(context, AggregatorID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range params.WasteTypes {
		gen.REPO.CreateAggregatorWasteType(context, gen.CreateAggregatorWasteTypeParams{
			AggregatorID: AggregatorID,
			WasteID:      v.ID,
			AlertLevel:   null.FloatFrom(v.AlertLevel).NullFloat64,
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Waste types set for aggregator",
	})
}

func (controller AggregatorController) GetWasteTypes(context *gin.Context) {

	parentID := context.Query("p")
	parentWasteTypeFilter_, _ := strconv.ParseUint(parentID, 10, 32)

	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	type Result struct {
		ID         int32          `json:"id"`
		Name       string         `json:"name"`
		IsActive   bool           `json:"is_active"`
		ParentID   sql.NullInt32  `json:"parent_id"`
		AlertLevel float64        `json:"alert_level"`
		CreatedAt  time.Time      `json:"created_at"`
		FilePath   sql.NullString `json:"file_path"`
		Parent     *Result        `json:"parent"`
	}

	type GetAllWasteTypesRow struct {
		ID         int32          `json:"id"`
		Name       string         `json:"name"`
		IsActive   bool           `json:"is_active"`
		IsSelected bool           `json:"is_selected"`
		AlertLevel float64        `json:"alert_level"`
		ParentID   sql.NullInt32  `json:"parent_id"`
		Parent     *Result        `json:"parent"`
		CreatedAt  time.Time      `json:"created_at"`
		FilePath   sql.NullString `json:"file_path"`
	}

	var results []GetAllWasteTypesRow

	if parentWasteTypeFilter_ <= 0 {
		wasteTypes, err := gen.REPO.GetAllWasteTypes(context)
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		var parentsFilter []Result

		for _, v := range wasteTypes {
			result := Result{}
			result.Name = v.Name
			result.ID = v.ID
			result.IsActive = v.IsActive
			result.ParentID = v.ParentID
			result.CreatedAt = v.CreatedAt
			result.FilePath = v.FilePath

			fmt.Print(v.ParentID)
			if v.ParentID.Valid {
				one, err := gen.REPO.GetOneWasteType(context, v.ParentID.Int32)
				if err == nil {
					result_ := Result{}
					result_.ID = one.ID
					result_.Name = one.Name
					result_.IsActive = one.IsActive
					result_.ParentID = one.ParentID
					result_.CreatedAt = one.CreatedAt
					result_.FilePath = one.FilePath

					result.Parent = &result_
				}
			} else {
				result.Parent = nil
			}
			parentsFilter = append(parentsFilter, result)
		}

		aggregatorWasteTypes, _ := gen.REPO.GetAggregatorWasteTypes(context, int32(auth.UserCompanyId.Int64))
		for _, v := range parentsFilter {

			isSelected := false
			item := GetAllWasteTypesRow{}
			for _, x := range aggregatorWasteTypes {
				if v.ID == x.WasteID {
					isSelected = true
					item.AlertLevel = x.AlertLevel.Float64
				}
			}

			item.CreatedAt = v.CreatedAt
			item.FilePath = v.FilePath
			item.ID = v.ID
			item.IsActive = v.IsActive
			item.Parent = v.Parent
			item.Name = v.Name
			item.ParentID = v.ParentID
			item.IsSelected = isSelected
			results = append(results, item)
		}

		context.JSON(http.StatusOK, gin.H{
			"error":       false,
			"waste_types": results,
		})
	} else {
		wasteTypes, err := gen.REPO.GetChildrenWasteTypes(context, sql.NullInt32{Int32: int32(parentWasteTypeFilter_), Valid: true})
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		aggregatorWasteTypes, _ := gen.REPO.GetAggregatorWasteTypes(context, int32(auth.UserCompanyId.Int64))

		for _, v := range wasteTypes {
			isSelected := false
			item := GetAllWasteTypesRow{}

			for _, x := range aggregatorWasteTypes {
				if v.ID == x.WasteID {
					isSelected = true
					item.AlertLevel = x.AlertLevel.Float64
				}
			}
			item.CreatedAt = v.CreatedAt
			item.FilePath = v.FilePath
			item.ID = v.ID
			item.IsActive = v.IsActive
			item.Name = v.Name
			item.ParentID = v.ParentID
			item.IsSelected = isSelected
			results = append(results, item)
		}

		context.JSON(http.StatusOK, gin.H{
			"error":       false,
			"waste_types": results,
		})
	}

}

func (controller AggregatorController) InsertAggregator(context *gin.Context) {
	var params CreateAggregatorParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := GetEmailUser(params.Email)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	context.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error":   true,
	// 		"message": "Error adding organization",
	// 	})
	// 	return
	// }
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email already exists`",
		})
		return
	}

	count, err := gen.REPO.GetDuplicateCompanies(context, gen.GetDuplicateCompaniesParams{
		Name:           strings.ToLower(params.Name),
		OrganizationID: sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another company has the same name withing the organization",
		})
		return
	}
	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}

	company, insertError := gen.REPO.InsertCompany(context, gen.InsertCompanyParams{
		Name:                         params.Name,
		Location:                     null.StringFrom(params.Location.Location).NullString,
		CountryID:                    country.ID,
		OrganizationID:               sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		IsActive:                     *params.IsActive,
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,

		CompanyType: 2, //params.Companytype,
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add aggregator",
		})
		return
	}

	_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,provider,role_id,user_organization_id,user_type,created_at,updated_at,password,user_company_id,is_company_super_admin) VALUES (:email,:first_name,:last_name,:provider,:role_id,:user_organization_id,:user_type,:created_at,:updated_at,:password,:user_company_id,:is_company_super_admin,:confirmed_at)`,
		map[string]interface{}{
			"email":                  params.Email,
			"first_name":             params.FirstName,
			"last_name":              params.LastName,
			"provider":               "email",
			"role_id":                3,
			"user_organization_id":   params.OrganizationID,
			"is_company_super_admin": true,
			"user_type":              9,
			"user_company_id":        company.ID,
			"password":               helpers.Functions{}.HashPassword(params.Password),
			"created_at":             time.Now(),
			"updated_at":             time.Now(),
			"confirmed_at":           time.Now(),
		})

	if params.LogoPath != "" {
		UploadController{}.SaveToUploadsTable(params.LogoPath, "companies", company.ID)
	}

	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Created Aggregator",
		"company": company, // Include the company details in the response
	})
}

func (controller AggregatorController) GetAllAggregators(context *gin.Context) {
	companies, err := gen.REPO.GetAllAggregators(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": companies,
	})
}

func (controller AggregatorController) GetAggregator(context *gin.Context) {
	id := context.Param("id")

	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	company, err := gen.REPO.GetCompany(context, int32(id_))

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"company": company,
	})
}

func (c AggregatorController) DeleteAggregator(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	err := gen.REPO.DeleteCompany(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Company successfully deleted",
	})
}

func (aggregatorController AggregatorController) UpdateAggregatorStatus(context *gin.Context) {
	var params UpdateAggregatorStatusParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateCompanyStatus(context, gen.UpdateCompanyStatusParams{
		IsActive: *params.IsActive,
		ID:       int32(params.ID),
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
		"message": "Successfully updated Aggregator",
		"status":  params.IsActive, // Use the variable for the parsed status
	})
}

func (aggregatorController AggregatorController) UpdateAggregator(context *gin.Context) {
	var params UpdateAggregatorDataParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	count, err := gen.REPO.GetDuplicateCompaniesWithoutID(context, gen.GetDuplicateCompaniesWithoutIDParams{
		Name:           strings.ToLower(params.Name),
		ID:             int32(params.ID),
		OrganizationID: sql.NullInt32{Int32: params.OrganizationID, Valid: true},
	})
	if len(count) > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another company has the same name withing the organization",
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

	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}

	// AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
	// 	Name:                         params.Name,
	// 	Location:                     null.StringFrom(params.Location.Location).NullString,
	// 	IsActive:                     *params.IsActive,
	// 	CountryID:                    country.ID,
	// 	OrganizationID:               sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
	// 	Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
	// 	Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
	// 	Region:                       null.StringFrom(params.Region).NullString,
	// 	CompanyType:                  2,

	// Update company status
	updateError := gen.REPO.UpdateCompany(context, gen.UpdateCompanyParams{
		IsActive:       *params.IsActive,
		Name:           params.Name,
		Location:       null.StringFrom(params.Location.Location).NullString,
		Region:         null.StringFrom(params.Region).NullString,
		Lat:            null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:            null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		CountryID:      country.ID,
		OrganizationID: sql.NullInt32{Int32: params.OrganizationID, Valid: params.OrganizationID != 0},
		CompanyType:    2,
		ID:             int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	if params.Password == "" {
		gen.REPO.UpdateUserFirstNameLastNameEmailRoleAndUserType(context, gen.UpdateUserFirstNameLastNameEmailRoleAndUserTypeParams{
			Email:     null.StringFrom(params.Email).NullString,
			RoleID:    sql.NullInt32{Int32: int32(3), Valid: true},
			UserType:  sql.NullInt16{Int16: int16(9), Valid: true},
			FirstName: sql.NullString{String: params.FirstName, Valid: true},
			LastName:  sql.NullString{String: params.LastName, Valid: true},
			ID:        params.UserID,
		})

	} else {
		gen.REPO.UpdateUserFirstNameLastNameEmailRoleUserTypeAndPassword(context, gen.UpdateUserFirstNameLastNameEmailRoleUserTypeAndPasswordParams{
			Email:     null.StringFrom(params.Email).NullString,
			RoleID:    sql.NullInt32{Int32: int32(3), Valid: true},
			UserType:  sql.NullInt16{Int16: int16(9), Valid: true},
			ID:        params.UserID,
			FirstName: sql.NullString{String: params.FirstName, Valid: true},
			LastName:  sql.NullString{String: params.LastName, Valid: true},
			Password:  null.StringFrom(helpers.Functions{}.HashPassword(params.Password)).NullString,
		})
	}

	UploadController{}.SaveToUploadsTable(params.LogoPath, "companies", int32(params.ID))

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated Company Data",
	})
}

func (aggregatorController AggregatorController) AddBuyer(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type Param struct {
		FirstName string `json:"first_name"  binding:"required"`
		LastName  string `json:"last_name"  binding:"required"`
		Company   string `json:"company"  binding:"required"`
		//CompanyID   int32           `json:"company_id"  binding:"required"`
		Location    models.Location `json:"location" binding:"required"`
		CallingCode string          `json:"calling_code" binding:"required"`
		Phone       string          `json:"phone" binding:"required"`
		Region      string          `json:"region"`
		IsActive    *bool           `json:"is_active" binding:"required"`
	}
	var params Param
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var duplicateBuyerCount int
	err = gen.REPO.DB.Get(&duplicateBuyerCount, gen.REPO.DB.Rebind("SELECT count(*) FROM buyers where company_id=? and LOWER(company)=? and LOWER(first_name)=? and LOWER(last_name)=?"),
		auth.UserCompanyId.Int64, strings.ToLower(params.Company), strings.ToLower(params.FirstName), strings.ToLower(params.LastName))

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if duplicateBuyerCount > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Duplicate buyer",
		})
		return
	}

	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}

	buyer, err := gen.REPO.CreateBuyer(context, gen.CreateBuyerParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		CountryID:                    sql.NullInt32{Int32: country.ID, Valid: true},
		IsActive:                     *params.IsActive,
		Company:                      null.StringFrom(params.Company).NullString,
		CallingCode:                  null.StringFrom(params.CallingCode).NullString,
		Phone:                        null.StringFrom(params.Phone).NullString,
		Location:                     null.StringFrom(params.Location.Location).NullString,
		CompanyID:                    int32(auth.UserCompanyId.Int64),
		CreatedAt:                    time.Now(),
		FirstName:                    params.FirstName,
		LastName:                     params.LastName,
		UpdatedAt:                    time.Now(),
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
	})

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": buyer,
	})
}

func (aggregatorController AggregatorController) AddSupplier(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type Param struct {
		FirstName string `json:"first_name"  binding:"required"`
		LastName  string `json:"last_name"  binding:"required"`
		Company   string `json:"company"  binding:"required"`
		//CompanyID   int32           `json:"company_id"  binding:"required"`
		Location    models.Location `json:"location" binding:"required"`
		CallingCode string          `json:"calling_code" binding:"required"`
		Phone       string          `json:"phone" binding:"required"`
		Region      string          `json:"region"`
		IsActive    *bool           `json:"is_active" binding:"required"`
	}
	var params Param
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var duplicateBuyerCount int
	err = gen.REPO.DB.Get(&duplicateBuyerCount, gen.REPO.DB.Rebind("SELECT count(*) FROM suppliers where company_id=? and LOWER(company)=? and LOWER(first_name)=? and LOWER(last_name)=?"),
		auth.UserCompanyId.Int64, strings.ToLower(params.Company), strings.ToLower(params.FirstName), strings.ToLower(params.LastName))

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if duplicateBuyerCount > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Duplicate supplier",
		})
		return
	}
	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}

	supplier, err := gen.REPO.CreateSupplier(context, gen.CreateSupplierParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		IsActive:                     *params.IsActive,
		CountryID:                    sql.NullInt32{Int32: country.ID, Valid: true},
		Company:                      null.StringFrom(params.Company).NullString,
		CallingCode:                  null.StringFrom(params.CallingCode).NullString,
		Phone:                        null.StringFrom(params.Phone).NullString,
		Location:                     null.StringFrom(params.Location.Location).NullString,
		CompanyID:                    int32(auth.UserCompanyId.Int64),
		CreatedAt:                    time.Now(),
		FirstName:                    params.FirstName,
		LastName:                     params.LastName,
		UpdatedAt:                    time.Now(),
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
	})

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": supplier,
	})
}

func (aggregatorController AggregatorController) UpdateSupplier(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type Param struct {
		FirstName   string          `json:"first_name"  binding:"required"`
		LastName    string          `json:"last_name"  binding:"required"`
		Company     string          `json:"company"  binding:"required"`
		Location    models.Location `json:"location" binding:"required"`
		CallingCode string          `json:"calling_code" binding:"required"`
		Phone       string          `json:"phone" binding:"required"`
		Region      string          `json:"region"`
		SupplierID  int32           `json:"supplier_id" binding:"required"`
		IsActive    *bool           `json:"is_active" binding:"required"`
	}
	var params Param
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var duplicateBuyerCount int
	err = gen.REPO.DB.Get(&duplicateBuyerCount, gen.REPO.DB.Rebind("SELECT count(*) FROM suppliers where company_id=? and LOWER(company)=? and LOWER(first_name)=? and LOWER(last_name)=? and id!=?"),
		auth.UserCompanyId.Int64, strings.ToLower(params.Company), strings.ToLower(params.FirstName), strings.ToLower(params.LastName), params.SupplierID)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if duplicateBuyerCount > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Duplicate supplier",
		})
		return
	}

	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}

	err = gen.REPO.UpdateSupplier(context, gen.UpdateSupplierParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		IsActive:                     *params.IsActive,
		CountryID:                    sql.NullInt32{Int32: country.ID, Valid: true},
		Company:                      null.StringFrom(params.Company).NullString,
		CallingCode:                  null.StringFrom(params.CallingCode).NullString,
		Phone:                        null.StringFrom(params.Phone).NullString,
		Location:                     null.StringFrom(params.Location.Location).NullString,
		CompanyID:                    int32(auth.UserCompanyId.Int64),
		FirstName:                    params.FirstName,
		LastName:                     params.LastName,
		ID:                           params.SupplierID,
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
	})

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Supplier updated successfully",
	})
}

func (aggregatorController AggregatorController) UpdateBuyer(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type Param struct {
		FirstName   string          `json:"first_name"  binding:"required"`
		LastName    string          `json:"last_name"  binding:"required"`
		Company     string          `json:"company"  binding:"required"`
		Location    models.Location `json:"location" binding:"required"`
		CallingCode string          `json:"calling_code" binding:"required"`
		Phone       string          `json:"phone" binding:"required"`
		Region      string          `json:"region"`
		BuyerID     int32           `json:"buyer_id" binding:"required"`
		IsActive    *bool           `json:"is_active" binding:"required"`
	}
	var params Param
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var duplicateBuyerCount int
	err = gen.REPO.DB.Get(&duplicateBuyerCount, gen.REPO.DB.Rebind("SELECT count(*) FROM buyers where company_id=? and LOWER(company)=? and LOWER(first_name)=? and LOWER(last_name)=? and id!=?"),
		auth.UserCompanyId.Int64, strings.ToLower(params.Company), strings.ToLower(params.FirstName), strings.ToLower(params.LastName), params.BuyerID)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if duplicateBuyerCount > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Duplicate buyer",
		})
		return
	}

	country, err := gen.REPO.GetCountryByName(context, params.Location.Country)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting country",
		})
		return
	}
	err = gen.REPO.UpdateBuyer(context, gen.UpdateBuyerParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		IsActive:                     *params.IsActive,
		Company:                      null.StringFrom(params.Company).NullString,
		CountryID:                    sql.NullInt32{Int32: country.ID, Valid: true},
		CallingCode:                  null.StringFrom(params.CallingCode).NullString,
		Phone:                        null.StringFrom(params.Phone).NullString,
		Location:                     null.StringFrom(params.Location.Location).NullString,
		CompanyID:                    int32(auth.UserCompanyId.Int64),
		FirstName:                    params.FirstName,
		LastName:                     params.LastName,
		ID:                           params.BuyerID,
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
	})

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Buyer updated successfully",
	})

}

func (aggregatorController AggregatorController) DeleteBuyer(context *gin.Context) {
	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	err := gen.REPO.DeleteBuyer(context, id32)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Buyer deleted successfully",
	})
}

func (aggregatorController AggregatorController) DeleteSupplier(context *gin.Context) {
	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	err := gen.REPO.DeleteSupplier(context, id32)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Supplier deleted successfully",
	})
}

func (aggregatorController AggregatorController) GetBuyers(context *gin.Context) {
	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")

	searchQuery := ""
	companyQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (first_name ilike " + "'%" + search + "%'" + " or last_name ilike " + "'%" + search + "%'" + "" + " or company ilike " + "'%" + search + "%')"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)

		//limitOffset = " LIMIT " + itemsPerPage + " OFFSET " + page
	}
	if companyID == "" {
		companyQuery = fmt.Sprint(" and company_id=", auth.UserCompanyId.Int64)
	} else {
		companyQuery = " and company_id=" + companyID
	}
	query := `
		select 
		buyers.id,
		buyers.first_name,
		buyers.last_name,
		buyers.company,
		buyers.location,
		buyers.is_active,
		countries.name as country_name,
		buyers.lat,
		buyers.lng,
		buyers.administrative_level_1_location,
		buyers.calling_code,
		buyers.phone
		from buyers 
		left join countries on countries.id = buyers.country_id
		where created_at is not null
	 ` + searchQuery + companyQuery + " order by created_at " + limitOffset

	var totalCount = 0
	gen.REPO.DB.Get(&totalCount, "select count(*) from buyers where created_at is not null"+companyQuery)
	logger.Log("AggregatorController/GetBuyers", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"total_count": totalCount,
		"content":     results,
	})
}

func (aggregatorController AggregatorController) GetSuppliers(context *gin.Context) {
	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")

	searchQuery := ""
	companyQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (first_name ilike " + "'%" + search + "%'" + " or last_name ilike " + "'%" + search + "%'" + "" + " or company ilike " + "'%" + search + "%')"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if companyID == "" {
		companyQuery = fmt.Sprint(" and  company_id=", auth.UserCompanyId.Int64)
	} else {
		companyQuery = " and company_id=" + companyID
	}
	query := `
		select 
		suppliers.id,
		suppliers.first_name,
		suppliers.last_name,
		suppliers.company,
		suppliers.location,
		suppliers.is_active,
		countries.name as country_name,
		suppliers.administrative_level_1_location,
		suppliers.calling_code,
		suppliers.phone
		from suppliers
		left join countries on countries.id = suppliers.country_id

		where created_at is not null
	 ` + searchQuery + companyQuery + " order by created_at " + limitOffset

	logger.Log("AggregatorController/GetSuppliers", query, logger.LOG_LEVEL_INFO)

	var totalCount = 0
	gen.REPO.DB.Get(&totalCount, "select count(*) from suppliers where created_at is not null"+companyQuery)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"total_count": totalCount,
		"content":     results,
	})
}

func (aggregatorController AggregatorController) MakeInventoryAdjustments(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	type WasteItem struct {
		ID             int32   `json:"id"  binding:"required"`
		Adjustment     float64 `json:"adjustment"  binding:"dive"`
		AdjustmentType string  `json:"adjustment_type"  binding:"required"`
	}
	type Param struct {
		// CompanyID  int32       `json:"company_id"  binding:"required"`
		Reason     string      `json:"reason"`
		WasteItems []WasteItem `json:"waste_items"  binding:"required"`
	}
	var params Param
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	var inventoryErrors []string
	for _, v := range params.WasteItems {
		if v.AdjustmentType == "negative" {
			wasteItem, err := gen.REPO.GetOneWasteType(context, v.ID)
			if err == nil {
				inventoryCount, _ := gen.REPO.InventoryItemCount(context, gen.InventoryItemCountParams{
					WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
					CompanyID:   int32(auth.UserCompanyId.Int64),
				})
				if inventoryCount == 0 {
					inventoryErrors = append(inventoryErrors, fmt.Sprint("Waste item ", wasteItem.Name, " does not exists in inventory"))
				} else {
					//insert
					item, _ := gen.REPO.GetInventoryItem(context, gen.GetInventoryItemParams{
						WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
						CompanyID:   int32(auth.UserCompanyId.Int64),
					})
					currentQuantity := item.TotalWeight
					if currentQuantity-v.Adjustment < 0 {
						inventoryErrors = append(inventoryErrors, fmt.Sprint("Invalid new total weight for waste ", wasteItem.Name, " new total weight will be ", currentQuantity-v.Adjustment))
					}
				}
			} else {
				inventoryErrors = append(inventoryErrors, "Error getting waste type")
			}
		}
	}
	if len(inventoryErrors) != 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": strings.Join(inventoryErrors, "\n"),
		})
		return
	}
	for _, v := range params.WasteItems {
		if err == nil {
			inventoryCount, _ := gen.REPO.InventoryItemCount(context, gen.InventoryItemCountParams{
				WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64),
			})
			if inventoryCount == 0 {
				//do insert
				err = gen.REPO.InsertToInventory(context, gen.InsertToInventoryParams{
					WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
					CompanyID:   int32(auth.UserCompanyId.Int64),
					TotalWeight: v.Adjustment,
				})
				if err != nil {
					logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
				} else {
					err = gen.REPO.InsertToInventoryAdjustments(context, gen.InsertToInventoryAdjustmentsParams{
						AdjustedBy:           int32(auth.ID.Int64),
						CompanyID:            int32(auth.UserCompanyId.Int64),
						AdjustmentAmount:     v.Adjustment,
						WasteTypeID:          v.ID,
						IsPositiveAdjustment: true,
						Reason:               null.StringFrom(params.Reason).NullString,
					})
					if err != nil {
						logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
					}
				}
			} else {
				//insert
				item, _ := gen.REPO.GetInventoryItem(context, gen.GetInventoryItemParams{
					WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
					CompanyID:   int32(auth.UserCompanyId.Int64),
				})

				currentQuantity := item.TotalWeight
				if v.AdjustmentType == "negative" {
					remainingWeight := currentQuantity - v.Adjustment
					logger.Log("Aggregator/MakeAdjustment", fmt.Sprint("[Negative adjustment] remaining weight ", remainingWeight), logger.LOG_LEVEL_INFO)

					err = gen.REPO.UpdateInventoryItem(context, gen.UpdateInventoryItemParams{
						TotalWeight: remainingWeight,
						ID:          item.ID,
					})

					if err != nil {
						logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
					} else {
						err = gen.REPO.InsertToInventoryAdjustments(context, gen.InsertToInventoryAdjustmentsParams{
							AdjustedBy:           int32(auth.ID.Int64),
							CompanyID:            int32(auth.UserCompanyId.Int64),
							AdjustmentAmount:    v.Adjustment,
							WasteTypeID:          v.ID,
							IsPositiveAdjustment: false,
							Reason:               null.StringFrom(params.Reason).NullString,
						})
						if err != nil {
							logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
						}
					}
				} else if v.AdjustmentType == "positive" {
					remainingWeight := currentQuantity + v.Adjustment
					logger.Log("Aggregator/MakeAdjustment", fmt.Sprint("[positive adjustment] remaining weight ", remainingWeight, "-", v.ID), logger.LOG_LEVEL_INFO)
					err = gen.REPO.UpdateInventoryItem(context, gen.UpdateInventoryItemParams{
						TotalWeight: remainingWeight,
						ID:          item.ID,
					})
					if err != nil {
						logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
					} else {
						err = gen.REPO.InsertToInventoryAdjustments(context, gen.InsertToInventoryAdjustmentsParams{
							AdjustedBy:           int32(auth.ID.Int64),
							CompanyID:            int32(auth.UserCompanyId.Int64),
							AdjustmentAmount:     v.Adjustment,
							IsPositiveAdjustment: true,
							WasteTypeID:          v.ID,
							Reason:               null.StringFrom(params.Reason).NullString,
						})
						if err != nil {
							logger.Log("Aggregator/MakeAdjustment", err.Error(), logger.LOG_LEVEL_ERROR)
						}
					}
				}
			}
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Inventory adjusted successfully",
	})
}

type PurchaseWasteItem struct {
	ID        int32   `json:"id"  binding:"required"`
	Weight    float64 `json:"weight"  binding:"required"`
	CostPerKG float64 `json:"cost_per_kg" binding:"required"`
	Amount    float64 `json:"amount" binding:"amount"`
}
type PurchaseWasteParam struct {
	SupplierID          int32               `json:"supplier_id"  binding:"required"`
	Date                string              `json:"date"  binding:"required"`
	PurchaseTotalAmount float64             `json:"purchase_total_amount"  binding:"required"` //Allow Partial Payments?
	PurchaseItems       []PurchaseWasteItem `json:"waste_items"  binding:"required"`
	PaymentMethod       int32               `json:"payment_method"  binding:"required"`
}

func (aggregatorController AggregatorController) PurchaseWasteFromSupplier(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	var params PurchaseWasteParam
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	date, parseDateError := time.Parse("2006-01-02", params.Date)
	if parseDateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": parseDateError.Error(),
		})
		return
	}
	if params.PaymentMethod == 1 {
		err := PurchaseWasteFromSupplierCash(params, auth, sql.NullTime{Time: date, Valid: true})
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Waste purchased successfully",
		})
	}
}

// for cash payments we automatically add items to the inventory since we assume that the user has paid full amount
func PurchaseWasteFromSupplierCash(param PurchaseWasteParam, auth *models.User, date sql.NullTime) error {
	var totalAmount = 0.0
	var totalWeight = 0.0
	ref := helpers.Functions{}.GetRandString(6)
	for _, v := range param.PurchaseItems {
		totalAmount += v.Amount // v.CostPerKG * v.Weight
		totalWeight += v.Weight
	}

	//create sell
	purchase, err := gen.REPO.CreatePurchase(context.Background(), gen.CreatePurchaseParams{
		Ref:         ref,
		CompanyID:   int32(auth.UserCompanyId.Int64),
		SupplierID:  param.SupplierID,
		TotalWeight: null.FloatFrom(totalWeight).NullFloat64,
		TotalAmount: null.FloatFrom(totalAmount).NullFloat64,
	})
	if err != nil {
		return err
	}
	//create sell items
	var errorSavingPurchaseItem = false
	for _, v := range param.PurchaseItems {
		_, err := gen.REPO.CreatePurchaseItem(context.Background(), gen.CreatePurchaseItemParams{
			CompanyID:   int32(auth.UserCompanyId.Int64),
			PurchaseID:  purchase.ID,
			WasteTypeID: v.ID,
			Weight:      null.FloatFrom(v.Weight).NullFloat64,
			CostPerKg:   null.FloatFrom(v.CostPerKG).NullFloat64,
			TotalAmount: null.FloatFrom(v.Weight * v.CostPerKG).Float64,
		})
		if err != nil {
			errorSavingPurchaseItem = true
		}
	}
	if errorSavingPurchaseItem {
		gen.REPO.DeletePurchase(context.Background(), purchase.ID)
		return errors.New("Error occured")
	}
	_, err = gen.REPO.MakePurchaseCashPayment(context.Background(), gen.MakePurchaseCashPaymentParams{
		Ref:             helpers.Functions{}.GetRandString(6),
		PurchaseID:      purchase.ID,
		PaymentMethod:   "CASH",
		Amount:          fmt.Sprint(totalAmount),
		CompanyID:       int32(auth.UserCompanyId.Int64),
		TransactionDate: date,
	})
	if err != nil {
		gen.REPO.DeletePurchase(context.Background(), purchase.ID)
		return err
	}

	//var errorSavingInventory = false
	for _, v := range param.PurchaseItems {
		item, err := gen.REPO.GetInventoryItem(
			context.Background(), gen.GetInventoryItemParams{
				WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64)})

		if err != nil && err == sql.ErrNoRows {
			gen.REPO.InsertToInventory(context.Background(), gen.InsertToInventoryParams{
				TotalWeight: v.Weight,
				CompanyID:   int32(auth.UserCompanyId.Int64),
				WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
			})
		} else if err != nil && err != sql.ErrNoRows {
			//errorSavingInventory = true
			logger.Log("AggregatorController/PurchaseWasteFromSupplierCash", fmt.Sprint("Error saving to inventory :: ", err.Error()), logger.LOG_LEVEL_ERROR)
		} else {
			currentQuantity := item.TotalWeight

			var remainingWeight = currentQuantity + v.Weight
			//update with the remaining weight
			gen.REPO.UpdateInventoryItem(context.Background(), gen.UpdateInventoryItemParams{
				TotalWeight: remainingWeight,
				ID:          item.ID,
			})
		}
	}

	return nil
}

type SellWasteItem struct {
	ID        int32   `json:"id"  binding:"required"`
	Weight    float64 `json:"weight"  binding:"required"`
	CostPerKG float64 `json:"cost_per_kg" binding:"required"`
	Amount    float64 `json:"amount" binding:"amount"`
}
type SellWasteParam struct {
	BuyerID         int32           `json:"buyer_id"  binding:"required"`
	Date            string          `json:"date"  binding:"required"`
	SellTotalAmount float64         `json:"sell_total_amount"  binding:"required"` //Allow Partial Payments?
	WasteItems      []SellWasteItem `json:"waste_items"  binding:"required"`
	PaymentMethod   int32           `json:"payment_method"  binding:"required"`
}

func (aggregatorController AggregatorController) SellWasteToBuyer(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	var params SellWasteParam
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	date, parseDateError := time.Parse("2006-01-02", params.Date)
	if parseDateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": parseDateError.Error(),
		})
		return
	}
	if params.PaymentMethod == 1 {
		err := SellWasteToBuyerCash(params, auth, sql.NullTime{Time: date, Valid: true})
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Waste sold successfully",
		})
	}
}

func (controller AggregatorController) SetSupplierActiveInActiveStatus(context *gin.Context) {
	type Param struct {
		ID       int32 `json:"id" binding:"required"`
		IsActive *bool `json:"is_active" binding:"required"`
	}
	var param Param
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err = gen.REPO.SetSupplierActiveInactiveStatus(context, gen.SetSupplierActiveInactiveStatusParams{
		IsActive: *param.IsActive,
		ID:       int32(param.ID),
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error setting supplier status",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Updated supplier status",
		})
	}
}

func (controller AggregatorController) SetBuyerActiveInActiveStatus(context *gin.Context) {
	type Param struct {
		ID       int32 `json:"id" binding:"required"`
		IsActive *bool `json:"is_active" binding:"required"`
	}
	var param Param
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err = gen.REPO.SetBuyerActiveInactiveStatus(context, gen.SetBuyerActiveInactiveStatusParams{
		IsActive: *param.IsActive,
		ID:       int32(param.ID),
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error setting buyer status",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Updated buyer status",
		})
	}
}

// for cash payments we automatically add items to the inventory since we assume that the user has paid full amount
func SellWasteToBuyerCash(param SellWasteParam, auth *models.User, date sql.NullTime) error {
	var totalAmount = 0.0
	var totalWeight = 0.0
	ref := helpers.Functions{}.GetRandString(6)
	for _, v := range param.WasteItems {
		totalAmount += v.CostPerKG * v.Weight
		totalWeight += v.Weight
	}
	var inventoryErrors []string

	for _, v := range param.WasteItems {
		wasteItem, err := gen.REPO.GetOneWasteType(context.Background(), v.ID)
		if err == nil {
			inventoryCount, _ := gen.REPO.InventoryItemCount(context.Background(), gen.InventoryItemCountParams{
				WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64),
			})
			if inventoryCount == 0 {
				inventoryErrors = append(inventoryErrors, fmt.Sprint("Waste item ", wasteItem.Name, " does not exists in inventory"))
			} else {
				//insert
				item, _ := gen.REPO.GetInventoryItem(context.Background(), gen.GetInventoryItemParams{
					WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
					CompanyID:   int32(auth.UserCompanyId.Int64),
				})
				currentQuantity := item.TotalWeight // strconv.ParseFloat(strings.TrimSpace(item.TotalWeight), 64)
				if currentQuantity-v.Weight < 0 {
					inventoryErrors = append(inventoryErrors, fmt.Sprint("Not enough items in the inventory for waste item ", wasteItem.Name, " current quantity is ", currentQuantity, " Kgs requested quantity is ", totalWeight, " kgs"))
				}
			}
		} else {
			inventoryErrors = append(inventoryErrors, "One of the waste types does not exist in the inventory")
		}
	}

	if len(inventoryErrors) != 0 {
		return errors.New(strings.Join(inventoryErrors, "\n"))
	}

	//create sell
	sale, err := gen.REPO.CreateSale(context.Background(), gen.CreateSaleParams{
		Ref:         ref,
		CompanyID:   int32(auth.UserCompanyId.Int64),
		BuyerID:     param.BuyerID,
		TotalWeight: null.FloatFrom(totalWeight).NullFloat64,
		TotalAmount: null.FloatFrom(totalAmount).NullFloat64,
	})
	if err != nil {
		return err
	}
	//create sell items
	var errorSavingSaleItem = false
	for _, v := range param.WasteItems {
		_, err := gen.REPO.CreateSaleItem(context.Background(), gen.CreateSaleItemParams{
			CompanyID:   int32(auth.UserCompanyId.Int64),
			SaleID:      sale.ID,
			WasteTypeID: v.ID,
			Weight:      null.FloatFrom(v.Weight).NullFloat64,
			CostPerKg:   null.FloatFrom(v.CostPerKG).NullFloat64,
			TotalAmount: v.Weight * v.CostPerKG,
		})
		if err != nil {
			errorSavingSaleItem = true
		}
	}
	if errorSavingSaleItem {
		gen.REPO.DeleteSale(context.Background(), sale.ID)
		return errors.New("Error occured")
	}
	_, err = gen.REPO.MakeCashPayment(context.Background(), gen.MakeCashPaymentParams{
		Ref:             helpers.Functions{}.GetRandString(6),
		SaleID:          sale.ID,
		PaymentMethod:   "CASH",
		Amount:          fmt.Sprint(totalAmount),
		CompanyID:       int32(auth.UserCompanyId.Int64),
		TransactionDate: date,
	})
	if err != nil {
		gen.REPO.DeleteSale(context.Background(), sale.ID)
		return err
	}

	for _, v := range param.WasteItems {
		item, _ := gen.REPO.GetInventoryItem(
			context.Background(), gen.GetInventoryItemParams{
				WasteTypeID: sql.NullInt32{Int32: v.ID, Valid: true},
				CompanyID:   int32(auth.UserCompanyId.Int64)})

		currentQuantity := item.TotalWeight

		var remainingWeight = currentQuantity - v.Weight
		//update with the remaining weight
		gen.REPO.UpdateInventoryItem(context.Background(), gen.UpdateInventoryItemParams{
			TotalWeight: remainingWeight,
			ID:          item.ID,
		})

	}

	return nil
}

func SellWasteToBuyerCashless() {

}

func (controller AggregatorController) ViewInventory(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	companyID := context.Query("cid")
	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")

	if companyID == "" {
		companyID = fmt.Sprint(auth.UserCompanyId.Int64)
	}
	searchQuery := ""
	limitOffset := ""
	companyQuery := " and  q.company_id=" + companyID

	if search != "" {
		searchQuery = " and (q.company_name ilike " + "'%" + search + "%'" + " or q.waste_name ilike " + "'%" + search + "%'" + ")"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}

	query := `
	 select * from 
	 (
		select 
        inventory.id,
		inventory.company_id,
		inventory.waste_type_id,
		inventory.total_weight,
		waste_types.name as waste_name,
		companies.name as company_name,

		coalesce(aggregator_waste_types.alert_level, 0) as alert_level
		
		from inventory

		left join aggregator_waste_types on aggregator_waste_types.waste_id = inventory.waste_type_id and aggregator_waste_types.aggregator_id=inventory.company_id
		inner join waste_types on waste_types.id=inventory.waste_type_id
		inner join companies on companies.id = inventory.company_id
	 ) as q where q.id is not null` + searchQuery + companyQuery + " order by q.id desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from inventory where id is not null and company_id=", companyID))
	logger.Log("AggregatorController/GetInventory", query, logger.LOG_LEVEL_INFO)

	if err != nil {
		logger.Log("AggregatorController/GetInventory", err.Error(), logger.LOG_LEVEL_ERROR)
	}
	results, err := utils.Select(gen.REPO.DB, query)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"total_count": totalCount,
		"content":     results,
	})
}

func (aggregatorController AggregatorController) GetSales(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")
	dateRangeStart := context.Query("sd")
	dateRangeEnd := context.Query("ed")

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.first_name ilike " + "'%" + search + "%'" + " or q.company_name ilike " + "'%" + search + "%'" + " or q.last_name ilike " + "'%" + search + "%'" + " or q.ref ilike " + "'%" + search + "%'" + ")"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if companyID == "" {
		companyQuery = fmt.Sprint(" and  q.company_id=", auth.UserCompanyId.Int64)
	} else {
		companyQuery = " and  q.company_id=" + companyID
	}
	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.sale_date as date)>='" + dateRangeStart + "' and cast(q.sale_date as date)<='" + dateRangeEnd + "'"
	}
	query := `
	 select * from 
	 (
		select 
		sales.id,
		sales.ref,
		sales.company_id,
		sales.buyer_id,
		sales.total_weight,
		sales.total_amount,
		sales.date as sale_date,
		buyers.first_name,
		buyers.last_name,
		buyers.company as buyer_company,
		companies.name as company_name
		from sales 

		inner join buyers on buyers.id=sales.buyer_id
		inner join companies on companies.id = sales.company_id
	 ) as q where q.sale_date is not null` + dateRangeQuery + searchQuery + companyQuery + " order by q.sale_date desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from sales where date is not null and company_id=", auth.UserCompanyId.Int64))

	//fmt.Println(err.Error())
	logger.Log("AggregatorController/GetSales", query, logger.LOG_LEVEL_INFO)

	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, v := range results {
		id, _ := v["id"]
		//fmt.Println(id)
		items, _ := utils.Select(gen.REPO.DB, fmt.Sprint("select *,waste_types.name as waste_name from sale_items join waste_types on waste_types.id=sale_items.waste_type_id where sale_items.sale_id=", id))
		//fmt.Println(items)
		v["items"] = items
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}

func (aggregatorController AggregatorController) GetUsers(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")

	searchQuery := ""
	companyQuery := ""

	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.first_name ilike " + "'%" + search + "%'" + " or q.email ilike " + "'%" + search + "%'" + " or q.last_name ilike " + "'%" + search + "%'" + " or q.role_name ilike " + "'%" + search + "%'" + ")"
	}

	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}
	if companyID == "" {
		companyID = fmt.Sprint(auth.UserCompanyId.Int64)
	}

	companyQuery = " and q.user_company_id=" + companyID

	query := `
	 select * from 
	 (
		select 
		users.id,
		users.first_name,
		users.last_name,
		users.user_company_id,
		users.email,
		users.is_active,
		users.calling_code,
		users.phone,
		users.role_id,
		users.created_at,
		roles.name as role_name

		from users

		inner join roles on users.role_id=roles.id

	 ) as q where q.created_at is not null and q.role_id!=3` + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from users where created_at is not null and users.role_id!=3 and user_company_id=", companyID))

	if err != nil {
		logger.Log("AggregatorController/GetUsers", err.Error(), logger.LOG_LEVEL_ERROR)
	}
	logger.Log("AggregatorController/GetUsers", query, logger.LOG_LEVEL_INFO)
	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}

func (aggregatorController AggregatorController) GetPurchases(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")
	dateRangeStart := context.Query("sd")
	dateRangeEnd := context.Query("ed")

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.first_name ilike " + "'%" + search + "%'" + " or q.company_name ilike " + "'%" + search + "%'" + " or q.last_name ilike " + "'%" + search + "%'" + " or q.ref ilike " + "'%" + search + "%'" + ")"
	}
	if itemsPerPage != "" && page != "" {
		limitOffset = " LIMIT " + itemsPerPage + " OFFSET " + page
	}
	if companyID == "" {
		companyQuery = fmt.Sprint(" and  q.company_id=", auth.UserCompanyId.Int64)
	} else {
		companyQuery = " and  q.company_id=" + companyID
	}
	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.purchase_date as date)>='" + dateRangeStart + "' and cast(q.purchase_date as date)<='" + dateRangeEnd + "'"
	}
	query := `
	 select * from 
	 (
		select 
		purchases.id,
		purchases.ref,
		purchases.company_id,
		purchases.supplier_id,
		purchases.total_weight,
		purchases.total_amount,
		purchases.date as purchase_date,
		suppliers.first_name,
		suppliers.last_name,
		suppliers.company as supplier_company,
		companies.name as company_name
		from purchases

		inner join suppliers on suppliers.id=purchases.supplier_id
		inner join companies on companies.id = purchases.company_id
	 ) as q where q.purchase_date is not null` + dateRangeQuery + searchQuery + companyQuery + " order by q.purchase_date desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from purchases where date is not null and company_id=", auth.UserCompanyId.Int64))

	logger.Log("AggregatorController/GetPurchases", query, logger.LOG_LEVEL_INFO)
	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	for _, v := range results {
		id, _ := v["id"]
		//fmt.Println(id)
		items, _ := utils.Select(gen.REPO.DB, fmt.Sprint("select *,waste_types.name as waste_name from purchase_items join waste_types on waste_types.id=purchase_items.waste_type_id where purchase_items.purchase_id=", id))
		//fmt.Println(items)
		v["items"] = items
	}
	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}

func (aggregatorController AggregatorController) GetInventoryAdjustments(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	companyID := context.Query("cid")
	dateRangeStart := context.Query("sd")
	dateRangeEnd := context.Query("ed")

	searchQuery := ""
	companyQuery := ""
	dateRangeQuery := ""
	limitOffset := ""

	if search != "" {
		searchQuery = " and (q.adjusted_by_first_name ilike " + "'%" + search + "%'" + " or q.adjusted_by_last_name ilike " + "'%" + search + "%'" + " or q.waste_name ilike " + "'%" + search + "%'" + " or q.reason ilike " + "'%" + search + "%'" + ")"
	}
	if itemsPerPage != "" && page != "" {
		itemsPerPage, _ := strconv.Atoi(context.Query("ipp"))
		page, _ := strconv.Atoi(context.Query("p"))

		offset := (page - 1) * itemsPerPage

		limitOffset = fmt.Sprint(" LIMIT ", itemsPerPage, " OFFSET ", offset)
	}

	if companyID == "" {
		companyID = fmt.Sprint(auth.UserCompanyId.Int64)
	}

	companyQuery = " and q.company_id=" + companyID

	if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeQuery = " and cast(q.created_at as date)>='" + dateRangeStart + "' and cast(q.created_at as date)<='" + dateRangeEnd + "'"
	}
	query := `
	 select * from 
	 (
		select 
		inventory_adjustments.id,
		inventory_adjustments.adjusted_by,
		inventory_adjustments.created_at,
		inventory_adjustments.company_id,
		inventory_adjustments.adjustment_amount,
		inventory_adjustments.is_positive_adjustment,
		inventory_adjustments.reason,
		inventory_adjustments.waste_type_id,

		users.first_name as adjusted_by_first_name,
		users.last_name as adjusted_by_last_name,
		companies.name as company_name,
		waste_types.name as waste_name

		from inventory_adjustments

		inner join waste_types on waste_types.id = inventory_adjustments.waste_type_id
		inner join users on  users.id=inventory_adjustments.adjusted_by
		inner join companies on companies.id = inventory_adjustments.company_id
	 ) as q where q.created_at is not null` + dateRangeQuery + searchQuery + companyQuery + " order by q.created_at desc " + limitOffset

	var totalCount = 0
	err := gen.REPO.DB.Get(&totalCount, fmt.Sprint("select count(*) from inventory_adjustments where created_at is not null and company_id=", companyID))

	logger.Log("AggregatorController/GetInventoryAdjustments", query, logger.LOG_LEVEL_INFO)
	results, err := utils.Select(gen.REPO.DB, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"content":     results,
		"total_count": totalCount,
	})
}
