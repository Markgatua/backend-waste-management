package seeder

import (
	"context"
	"fmt"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"
)

func Run() {
	appSettings, err := utils.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", "Error getting app settings", logger.LOG_LEVEL_ERROR)
		return
	}

	for _, v := range appSettings.Connections {
		fmt.Println("===================>",v.ConnectionString)
		conn, err := pgx.Connect(context.Background(), v.ConnectionString)
		if err != nil {
			logger.Log("SEEDER",fmt.Sprint("Unable to connect to database: %v",err),logger.LOG_LEVEL_ERROR)
		}else{
			q := queries.NewQuerier(conn)
			BloodGroupSeeder{}.Run(q)
			FamilyReltionSeeder{}.Run(q)
			ReligionSeeder{}.Run(q)
			GenderSeeder{}.Run(q)
			UserTitlesSeeder{}.Run(q)
			CountriesSeeder{}.Run(q)
			PermissionsSeeder{}.Run(q)
		}
		defer conn.Close(context.Background())
	}

}
