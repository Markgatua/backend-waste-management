package controllers

import (
	"database/sql"
	"net/http"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type RoleHasPermissions struct{}

type AssignPermissionParams struct {
	RoleID       int32 `json:"role_id"  binding:"required"`
	PermissionID int32 `json:"permission_id"  binding:"required"`
}
type RevokePermissionParams struct{
	RoleID       int32 `json:"role_id"  binding:"required"`
	PermissionID int32 `json:"permission_id"  binding:"required"`
}

func (roleHasPermissions RoleHasPermissions) AssignPermission(context *gin.Context) {
	var params AssignPermissionParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	insertError := gen.REPO.AssignPermission(context, gen.AssignPermissionParams{
		RoleID: sql.NullInt32{Int32: params.RoleID, Valid: true},
		PermissionID: sql.NullInt32{Int32: params.PermissionID, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Assign permission",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Assigned Permission",
	})

}

func (roleHasPermissions RoleHasPermissions) RevokePermission(context *gin.Context) {
	var params RevokePermissionParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	insertError := gen.REPO.RevokePermission(context, gen.RevokePermissionParams{
		RoleID: sql.NullInt32{Int32: params.RoleID, Valid: true},
		PermissionID: sql.NullInt32{Int32: params.PermissionID, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Revoke Permission",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Revoked Permission",
	})

}