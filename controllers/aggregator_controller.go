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
	"ttnmwastemanagementsystem/models"

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
		FirstName   string          `json:"first_name"  binding:"required"`
		LastName    string          `json:"last_name"  binding:"required"`
		Company     string          `json:"company"  binding:"required"`
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

	buyer, err := gen.REPO.CreateBuyer(context, gen.CreateBuyerParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
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
	err = gen.REPO.DB.Get(&duplicateBuyerCount, "SELECT count(*) FROM buyers where company_id=? and lower(comapny)=? and lower(FirstName)=? and lower(LastName)=? where id!=?",
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

	err = gen.REPO.UpdateBuyer(context, gen.UpdateBuyerParams{
		Lat:                          null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng:                          null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		Region:                       null.StringFrom(params.Region).NullString,
		IsActive:                     *params.IsActive,
		Company:                      null.StringFrom(params.Company).NullString,
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
		"message": "Buyer updated successfully",
	})
}

func (aggregatorController AggregatorController) GetBuyers(context *gin.Context) {
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
		"message": "Buyer updated successfully",
	})
}

type SellWasteItem struct {
	ID        int32   `json:"id"  binding:"required"`
	Weight    float64 `json:"weight"  binding:"required"`
	CostPerKG float64 `json:"cost_per_kg" binding:"required"`
	Amount    float64 `json:"amount" binding:"amount"`
}
type SellWasteParam struct {
	BuyerID         int32           `json:"buyer_id"  binding:"required"`
	Date            null.Time       `json:"date"  binding:"required"`
	SellTotalAmount float64         `json:"sell_total_amount"  binding:"sell_total_amount"` //Allow Partial Payments?
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
	if params.PaymentMethod == 1 {
		err := SellWasteToBuyerCash(params, auth)
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Error selling waste",
		})
	}
}

// for cash payments we automatically add items to the inventory since we assume that the user has paid full amount
func SellWasteToBuyerCash(param SellWasteParam, auth *models.User) error {
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
				currentQuantity, _ := strconv.ParseFloat(strings.TrimSpace(item.TotalWeight), 64)
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
		TotalWeight: null.StringFrom(fmt.Sprint(totalWeight)).NullString,
		TotalAmount: null.StringFrom(fmt.Sprint(totalAmount)).NullString,
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
			Weight:      null.StringFrom(fmt.Sprint(v.Weight)).NullString,
			CostPerKg:   null.StringFrom(fmt.Sprint(v.CostPerKG)).NullString,
			TotalAmount: fmt.Sprint(v.Weight * v.CostPerKG),
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
		CompanyID:       int32(auth.UserCompanyId.Int64),
		TransactionDate: param.Date.NullTime,
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

		currentQuantity, _ := strconv.ParseFloat(strings.TrimSpace(item.TotalWeight), 64)

		var remainingWeight = currentQuantity - v.Weight
		//update with the remaining weight
		gen.REPO.UpdateInventoryItem(context.Background(), gen.UpdateInventoryItemParams{
			TotalWeight: fmt.Sprint(remainingWeight),
			ID:          item.ID,
		})

	}

	return nil
}

func SellWasteToBuyerCashless() {

}

func (aggregatorController AggregatorController) GetSales(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	search := context.Query("s")
	itemsPerPage := context.Query("ipp")
	page := context.Query("p")
	//sortBy := context.Query("sort_by")
	//orderBy := context.Query("order_by")
	roleID := context.Query("role_id")

	searchQuery := ""
	if search != "" {
		searchQuery = " and WHERE to_tsvector(body) @@ to_tsquery('" + search + "')"
	}
	limitOffset := " LIMIT " + itemsPerPage + " OFFSET " + page
	roleIDQuery := ""
	if roleID != "" {
		roleIDQuery = " and where users.role_id=" + roleID
	}
	companyQuery := fmt.Sprint(" and where users.user_company_id=", auth.UserCompanyId.Int64)

	query := `select users.id, users.first_name, users.last_name, users.email, users.avatar_url, users.calling_code, users.phone, users.is_active, roles.name as role_name,
    roles.id as role_id from users inner join roles on users.role_id = roles.id where users.email not ilike 'superadmin@admin.com'
    and users.is_main_organization_user = false ` + companyQuery + roleIDQuery + searchQuery + limitOffset

}