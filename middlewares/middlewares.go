package middlewares

import (
	"fmt"
	"net/http"
	_"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		functions := helpers.Functions{}
		err := functions.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":true,
				"message":"Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}


func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		functions := helpers.Functions{}
		user,err:=functions.CurrentUserFromToken(c)
		if err!=nil{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":true,
				"message":"You do not have necessary permissions to perform the above action",
			})
			c.Abort()
			return
		}
		roleID:=user.RoleId;

		fmt.Println(roleID)
		// var permissions = [];

		// err = gen.REPO.DB.Get(&permissions, gen.REPO.DB.Rebind("select * from role_has_permissions where role_id=?"), permissions)
		// fmt.Println("*********************")
		// fmt.Println(err)

		c.Next()
	}
}
