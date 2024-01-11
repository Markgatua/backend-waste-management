package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type RolesController struct{}

type UpdateRoleParams struct {
	ID     		  int `json:"id"  binding:"required"`
	Name 		  string `json:"name"  binding:"required"`
	GuardName     string `json:"guard_name"  binding:"required"`
	Description   string `json:"description"  binding:"required"`
	DeletedAt	  bool   `json:"delete"`
}

func(rolesController  RolesController) GetRoles(context *gin.Context){
	roles, err := gen.REPO.GetRoles(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"Roles":roles,
	})
}

func(rolesController RolesController) GetRole(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	role, err := gen.REPO.GetRole(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"Role": role,
	})
}

func (rolesController RolesController) UpdateRole(context *gin.Context) {
	var params UpdateRoleParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	fmt.Println(params.DeletedAt)
	var isToDelete null.Time

	if (params.DeletedAt == true) {
		isToDelete = null.TimeFrom(time.Now())
	} else {
		
	}

	updateError := gen.REPO.UpdateRole(context, gen.UpdateRoleParams{
		Name: params.Name,
		ID:     int32(params.ID),
		GuardName: params.GuardName,
		Description: null.StringFrom(params.Description).NullString,
		DeletedAt: isToDelete.NullTime,
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
		"message": "Successfully updated Role",
	})
}