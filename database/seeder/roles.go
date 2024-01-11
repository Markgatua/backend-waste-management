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
				description := v["description"]
				guardName := v["guard_name"]
		

				q.InsertRole(context.Background(), gen.InsertRoleParams{
					Name:            fmt.Sprint(name),
					Description:     sql.NullString{String: fmt.Sprint(description), Valid: true},
					GuardName: 		 fmt.Sprint(guardName),

				})
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
