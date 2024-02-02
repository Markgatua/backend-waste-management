package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type UsersController struct{}

type UpdateUserParams struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	ID        uint64 `json:"id" binding:"required"`
}

type SetActiveInactiveStatusParam struct {
	UserID   int64 `json:"user_id" binding:"required"`
	IsActive *bool `json:"is_active" binding:"required"`
}

func (controller UsersController) SetActiveInActiveStatus(context *gin.Context) {
	var param SetActiveInactiveStatusParam
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetUserByID(param.UserID)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User does not exist",
		})
		return
	}
	err = gen.REPO.UpdateUserIsActive(context, gen.UpdateUserIsActiveParams{
		IsActive: sql.NullBool{Bool: *param.IsActive, Valid: true},
		ID:       int32(param.UserID),
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

func (usersController UsersController) GetAllMainOrganizationUsers(context *gin.Context) {
	users, err := gen.REPO.GetAllMainOrganizationUsers(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"users": users,
	})
}

func (usersController UsersController) GetAllUsers(context *gin.Context) {

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

	var results []gen.GetAllUsersRow

	err := gen.REPO.DB.Select(&results, query)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// users, err := gen.REPO.GetAllUsers(context)
	// if err != nil {
	// 	context.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error":   true,
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"users": results,
	})
}

func (usersController UsersController) GetMainOrganizationUser(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	user, err := gen.REPO.GetMainOrganizationUser(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
	})
}

func (usersController UsersController) GetUser(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	user, err := gen.REPO.GetMainOrganizationUser(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
	})
}

func (usersController UsersController) UpdateMainOrganizationUser(context *gin.Context) {
	var params UpdateUserParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	updateError := gen.REPO.UpdateMainOrganizationUser(context, gen.UpdateMainOrganizationUserParams{
		FirstName: null.StringFrom(params.FirstName).NullString,
		LastName:  null.StringFrom(params.LastName).NullString,
		ID:        int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Successfully updated user",
		})
		return
	}
}

func (usersController UsersController) GetUsersWithRole(context *gin.Context) {
	users, err := gen.REPO.GetUsersWithRole(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"users": users,
	})
}

func (usersController UsersController) GetCompanyUsers(context *gin.Context) {
	id := context.Param("id")

	id_, _ := strconv.ParseUint(id, 10, 32)

	users, err := gen.REPO.GetCompanyUsers(context, sql.NullInt32{Int32: int32(id_), Valid: true})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"users": users,
	})
}
