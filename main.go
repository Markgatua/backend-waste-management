package main

import (
	_ "fmt"
	"tutorial/controllers"
	"tutorial/gen"

	"github.com/gin-gonic/gin"
)
func main()  {
	router := gin.Default()
	gen.LoadRepo()

	usersController:=controllers.UsersController{}
	router.GET("/users",usersController.GetAllUsers)
	// router.POST("update/user",usersController.UpdateUSer)
	// router.GET("/users/roles",usersController.GetUsersWithRole)
	router.GET("/user/:id", usersController.GetUser)


	router.Run()
	//fmt.Println("Hello dabid")
}