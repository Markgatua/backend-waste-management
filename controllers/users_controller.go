package controllers

import (
	"net/http"
	"strconv"
	"ttnmwastemanagementsystem/gen"

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

func (usersController UsersController) GetAllTTNMUsers(context *gin.Context) {
	users, err := gen.REPO.GetAllTTNMUsers(context)
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

func (usersController UsersController) GetTTNMUser(context *gin.Context) {
	id := context.Param("id")

	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	user, err := gen.REPO.GetTTNMUser(context, int32(id_))

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

func (usersController UsersController) UpdateTTNMUser(context *gin.Context) {
	var params UpdateUserParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	updateError := gen.REPO.UpdateTTNMUser(context, gen.UpdateTTNMUserParams{
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
