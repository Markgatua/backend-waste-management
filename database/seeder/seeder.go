package seeder

import (
	"database/sql"
	"fmt"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"

	_ "github.com/lib/pq"
)

func Run() {
	appSettings, err := utils.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", "Error getting app settings", logger.LOG_LEVEL_ERROR)
		return
	}
	conn, err := sql.Open("postgres", appSettings.DBMasterConnectionString)
	defer conn.Close()
	if err != nil {
		logger.Log("SEEDER", fmt.Sprint("Unable to connect to database: %v", err), logger.LOG_LEVEL_ERROR)
	} else {
		queries := gen.New(conn)
		CountriesSeeder{}.Run(queries)
		PermissionsSeeder{}.Run(queries)
	}

}
