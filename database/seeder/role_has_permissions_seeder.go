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

type RoleHasPermissionsSeeder struct{}

func (roleHasPermissionsSeeder RoleHasPermissionsSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/ROLE HAS PERMISSIONS SEEDER]", "=======Seeding ROLE HAS PERMISSIONS======", logger.LOG_LEVEL_INFO)

	jsonFile, err := os.Open("assets/data/role_has_permissions.json")
	if err == nil {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]map[string]interface{}
		unmarshalError := json.Unmarshal(byteValue, &result)
		if unmarshalError == nil {
			for _, v := range result {
				
				roleId := int32( v["role_id"].(float64))
				permissionId := int32(v["permission_id"].(float64))

				count, err := q.GetDuplicateRoleHasPermission(context.Background(), gen.GetDuplicateRoleHasPermissionParams{
					RoleID:      sql.NullInt32{Int32: roleId, Valid: true},
					PermissionID:      sql.NullInt32{Int32: permissionId, Valid: true},
				})
				if err == nil {
					// var Role gen.Role
					if count == 0 {
						q.AssignPermission(context.Background(), gen.AssignPermissionParams{
							RoleID:      sql.NullInt32{Int32: roleId, Valid: true},
							PermissionID:      sql.NullInt32{Int32: permissionId, Valid: true},
						})
					}else{

					}
				}
					
				//fmt.Println(err.Error())
			}
		} else {
			fmt.Println(unmarshalError.Error())
		}
	} else {
		fmt.Println("Error reading from json file -- ", err.Error())
	}

	defer jsonFile.Close()
}