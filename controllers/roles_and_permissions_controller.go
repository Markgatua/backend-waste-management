package controllers

import (
	"context"
	"fmt"
	"net/http"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/utils"

	// "scms/database/gen"
	// "scms/repo"
	// "scms/src/controller/school"
	// "scms/src/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

type RoleAndPermissionsController struct{}

type AssignPermissionsToRoleParams struct {
	RoleID      int64   `json:"role_id" binding:"required"`
	Permissions []int32  `json:"permissions" binding:"required"`
}

type RemovePermissionsFromRoleParams struct {
	RoleID      int64   `json:"role_id" binding:"required"`
	Permissions []int32 `json:"permissions" binding:"required"`
}
type UpdateRoleParams struct {
	RoleID      int64  `json:"role_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	IsActive    *bool   `json:"is_active" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type AddRoleParams struct {
	Name        string `json:"name" binding:"required"`
	IsActive    *bool   `json:"is_active" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (c RoleAndPermissionsController) RemovePermissionsFromRole(context *gin.Context) {
	params := RemovePermissionsFromRoleParams{}
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = RemovePermissionsFromRole(int32(params.RoleID), params.Permissions)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Permissions removed for the specified roles",
	})
}
func (c RoleAndPermissionsController) GetAllPermissions(context *gin.Context) {
	permissions, err :=  gen.REPO.GetAllPermissions(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error getting permissions",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": permissions,
	})
}
func (c RoleAndPermissionsController) AssignPermissionsToRole(context *gin.Context) {
	assignPermissionsToRoleParam := AssignPermissionsToRoleParams{}
	err := context.ShouldBindJSON(&assignPermissionsToRoleParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = SetRolePermissions(assignPermissionsToRoleParam.RoleID, assignPermissionsToRoleParam.Permissions)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Permissions assigned for the specified roles",
	})
}

func (c RoleAndPermissionsController) GetRoles(context *gin.Context) {
	roles, err :=  gen.REPO.GetRoles(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   false,
			"message": "Error getting roles",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"content": roles,
		})
	}
}

func (c RoleAndPermissionsController) GetRole(context *gin.Context) {
	id, _ := context.Params.Get("id")

	var id32 int32
	fmt.Sscan(id, &id32)
	role, err :=  gen.REPO.GetRole(context, id32)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   false,
			"message": "Error getting role",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"content": role,
		})
	}
}

func (c RoleAndPermissionsController) DeleteRole(context *gin.Context) {
	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)
	if id32 == 1 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Cannot delete specified role",
		})
		return
	}
	err :=  gen.REPO.DeleteRole(context, id32)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error deleting role",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Role deleted successfully",
		})
	}
}

func (c RoleAndPermissionsController) UpdateRole(context *gin.Context) {
	updateRoleParams := UpdateRoleParams{}
	err := context.ShouldBindJSON(&updateRoleParams)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err =  gen.REPO.UpdateRole(context, gen.UpdateRoleParams{
		RoleID:      int32(updateRoleParams.RoleID),
		Description: null.StringFrom(updateRoleParams.Description).NullString,
		Name:        updateRoleParams.Name,
		IsActive:    *updateRoleParams.IsActive,
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error updating role",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Role updated successfully",
	})

}
func (c RoleAndPermissionsController) AddRole(context *gin.Context) {

	nextID := utils.GetNextTableID("roles")

	params := AddRoleParams{}
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	role, err :=  gen.REPO.CreateRole(context, gen.CreateRoleParams{
		Description: null.StringFrom(params.Description).NullString,
		Name:        params.Name,
		IsActive:    *params.IsActive,
		ID:          nextID,
	})
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error creating role",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": role,
	})
}

func (c RoleAndPermissionsController) GetRolePermissions(context *gin.Context) {
	id, _ := context.Params.Get("role_id")
	var id32 int32
	fmt.Sscan(id, &id32)
	permissions, err := GetPermissionsForRole(id32)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error deleting role",
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": permissions,
	})

}

// this function check if role exists or not
func RoleExists(roleID int64) (bool, error) {
	// count, err := schoolConnection.Queries.RoleExists(context.Background(), int32(roleID))
	// if err != nil {
	// 	return false, err
	// }
	// if count > 0 {
	// 	return true, nil
	// } else {
	// 	return false, nil
	// }
	return false, nil
}

func GetPermissionsForRole(roleId int32) ([]gen.GetPermissionsForRoleIDRow, error) {
	val, err :=  gen.REPO.GetPermissionsForRoleID(context.Background(), roleId)
	return val, err
}

func GetActionsFromPermissions(permissions []gen.GetPermissionsForRoleIDRow) []string {
	actions := []string{}
	for _, v := range permissions {
		actions = append(actions, v.Action)
	}
	return actions
}
func SetRolePermissions(roleID int64, permissions []int32) error {
	for _, v := range permissions {

		gen.REPO.AssignPermissionToRole(context.Background(), gen.AssignPermissionToRoleParams{
			RoleID:       int32(roleID),
			PermissionID: int32(v),
		})
		//fmt.Println(err.Error())
	}
	return nil
}

func RemovePermissionsFromRole(roleID int32, permission []int32) error {
	err :=  gen.REPO.RemovePermissionsFromRole(context.Background(), gen.RemovePermissionsFromRoleParams{
		PermissionIds: permission,
		RoleID:        roleID,
	})
	if err != nil {
		//fmt.Println(err.Error())
		return err
	}
	return nil
}