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

type RolesSeeder struct{}

func (rolesSeeder RolesSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/ROLES SEEDER]", "=======Seeding ROLES======", logger.LOG_LEVEL_INFO)

	jsonFile, err := os.Open("assets/data/roles.json")
	if err == nil {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]map[string]interface{}
		unmarshalError := json.Unmarshal(byteValue, &result)
		if unmarshalError == nil {
			for _, v := range result {
				name := v["name"]
				roleId := v["role_id"]
				description := v["description"]
				guardName := v["guard_name"]

				count, err := q.GetDuplicateRole(context.Background(), int32(roleId.(float64)))
				if err == nil {
					// var Role gen.Role
					if count == 0 {
						q.InsertRole(context.Background(), gen.InsertRoleParams{
							RoleID:      int32(roleId.(float64)),
							Name:        fmt.Sprint(name),
							Description: sql.NullString{String: fmt.Sprint(description), Valid: true},
							GuardName:   fmt.Sprint(guardName),
						})
					} else {
						q.UpdateRole(context.Background(), gen.UpdateRoleParams{
							RoleID:      int32(roleId.(float64)),
							Name:        fmt.Sprint(name),
							Description: sql.NullString{String: fmt.Sprint(description), Valid: true},
							GuardName:   fmt.Sprint(guardName),
						})

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
