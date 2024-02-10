package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type PermissionsSeeder struct{}

type Permissions []Permission

func UnmarshalPermissions(data []byte) (Permissions, error) {
	var r Permissions
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Permissions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Permission struct {
	Module      string              `json:"module"`
	Key         string              `json:"key"`
	Permissions []PermissionElement `json:"permissions"`
}

type PermissionElement struct {
	Name      string `json:"name"`
	SubModule string `json:"sub_module"`
	Action    string `json:"action"`
}

func (permissionsSeeder PermissionsSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/PERMISSIONS SEEDER]", "======= Seeding permissions======", logger.LOG_LEVEL_INFO)
	jsonFile, err := os.Open("assets/data/permissions.json")
	if err == nil {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		permissions, err := UnmarshalPermissions(byteValue)
		if err != nil {
			logger.Log("[SEEDER/PERMISSIONS SEEDER]", fmt.Sprint("Error reading permissions file"), logger.LOG_LEVEL_ERROR)
			return
		}
		permissionActions := []string{}
		for _, v := range permissions {
			for _, permissionElement := range v.Permissions {
				q.CreatePermission(context.Background(), gen.CreatePermissionParams{
					Name:      permissionElement.Name,
					Action:    permissionElement.Action,
					Module:    v.Module,
					Submodule: sql.NullString{String: permissionElement.SubModule, Valid: true},
				})
				permissionActions = append(permissionActions, permissionElement.Action)
			}
		}

		fmt.Println(permissionActions)
		//delete permissions which are not captured in the file
		permissionIDs := []int32{}
		q.DeletePermissionByActions(context.Background(), permissionActions)
		permissions_, err := q.GetAllPermissions(context.Background())
		if err == nil {
			for _, v := range permissions_ {
				permissionIDs = append(permissionIDs, v.ID)
			}
		}

		//superAdmin,_ := q.GetUserByEmail(context.Background(),sql.NullString{String: "superadmin@admin.com",Valid: true})
		//insert super admin permissions
		for _, v := range permissionIDs {
			q.AssignPermissionToRole(context.Background(), gen.AssignPermissionToRoleParams{
				RoleID:       3,
				PermissionID: v,
			})
			q.AssignPermissionToRole(context.Background(), gen.AssignPermissionToRoleParams{
				RoleID:       12,
				PermissionID: v,
			})
		}
	} else {
		fmt.Println("Error reading from json file -- ", err.Error())
	}
	defer jsonFile.Close()
}
