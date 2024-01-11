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
	PermissionId int32 `json:"permission_id"`
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
		for _, v := range permissions {

			for _, permissionElement := range v.Permissions {

				q.CreatePermission(context.Background(), gen.CreatePermissionParams{
					PermissionID: int32(permissionElement.PermissionId),
					Name:      permissionElement.Name,
					GuardName: permissionElement.Action,
					Module:    v.Module,
					Submodule: sql.NullString{String:fmt.Sprint(permissionElement.SubModule),Valid: true},
				})
			}

		}

	} else {
		fmt.Println("Error reading from json file -- ", err.Error())
	}

	defer jsonFile.Close()
}
