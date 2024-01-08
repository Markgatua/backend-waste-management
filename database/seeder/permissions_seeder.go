package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"scms/database/queries"
	"scms/src/logger"
	_ "scms/src/utils"

	"github.com/jackc/pgtype"
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

func (permissionsSeeder PermissionsSeeder) Run(q *queries.DBQuerier) {
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
				q.CreatePermission(context.Background(), queries.CreatePermissionParams{
					Name: pgtype.Varchar{String: permissionElement.Name,Status: pgtype.Present},
					GuardName: pgtype.Varchar{String: permissionElement.Action,Status: pgtype.Present},
					Module: pgtype.Varchar{String: v.Module,Status: pgtype.Present},
					Submodule: pgtype.Varchar{String: permissionElement.SubModule,Status: pgtype.Present},
				})
			}

		}

	} else {
		fmt.Println("Error reading from json file -- ", err.Error())
	}

	defer jsonFile.Close()
}
