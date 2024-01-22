package controllers

import (
	"database/sql"
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


func(usersController  UsersController) GetAllUsers(context *gin.Context){
	users, err := gen.REPO.GetAllUsers(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"users":users,
	})
}

func(usersController UsersController) GetUser(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);
	println("------------------------------",id_)
	user, err := gen.REPO.GetUser(context, int32(id_))

	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK,gin.H{
		"error":  false,
		"user": user,
	})
}

func(usersController UsersController) UpdateUSer(context *gin.Context){
	var params UpdateUserParams 
	err := context.ShouldBindJSON(&params)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":err.Error(),
		})
		return
	}

	updateError := gen.REPO.UpdateUser(context,gen.UpdateUserParams{
		FirstName: null.StringFrom(params.FirstName).NullString,
		LastName: null.StringFrom(params.LastName).NullString,
		ID: int32(params.ID),
	})
	if updateError!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
			"error":true,
			"message":updateError.Error(),
		})
		return
	}else{
		context.JSON(http.StatusOK,gin.H{
			"error":false,
			"message":"Successfully updated user",
		})
		return
	}
}



func(usersController  UsersController) GetUsersWithRole(context *gin.Context){
	users, err := gen.REPO.GetUsersWithRole(context)
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"users":users,
	})
}

func(usersController  UsersController) GetCompanyUsers(context *gin.Context){
	id :=  context.Param("id")

	id_,_ :=strconv.ParseUint(id,10,32);

	users, err := gen.REPO.GetCompanyUsers(context, sql.NullInt32{Int32: int32(id_),Valid: true});
	if err!=nil{
		context.JSON(http.StatusUnprocessableEntity,gin.H{
		   "error":true,
		   "message":err.Error(),	
		})
		return
	}
	
	context.JSON(http.StatusOK,gin.H{
		"error":false,
		"users":users,
	})
}