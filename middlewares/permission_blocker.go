package middlewares

import (
	"net/http"
	"ttnmwastemanagementsystem/controllers"
	"ttnmwastemanagementsystem/helpers"

	"slices"

	"github.com/gin-gonic/gin"
)

func PermissionBlockerMiddleware(action string) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth, err := helpers.Functions{}.CurrentUserFromToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Error getting user",
			})
		}
		permissions, err := controllers.GetPermissionsForRole(int32(auth.RoleId.Int64))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Error getting permissions",
			})
		}
		actionList :=  controllers.GetActionsFromPermissions(permissions)
		if !slices.Contains(actionList, action){
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "You do not have necessary permission to perform the action specified",
			})
		}

		c.Set("user", auth)
		c.Next()
	}
}