package seeder

import (
	"database/sql"
	"fmt"
	"ttnmwastemanagementsystem/appsettings"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"

	_ "github.com/lib/pq"
)

func Run() {
	appSettings, err := appsettings.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", fmt.Sprint("Error getting app settings::",err.Error()), logger.LOG_LEVEL_ERROR)
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
		TtnmOrganizationSeeder{}.Run(queries)
		RoleSeeder{}.Run(queries)
		RoleHasPermissionsSeeder{}.Run(queries)
	}

}
