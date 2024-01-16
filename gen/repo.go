package gen

import (
	"fmt"
	"log"
	"ttnmwastemanagementsystem/appsettings"
	"ttnmwastemanagementsystem/logger"

	// _"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var REPO *Repo

type Repo struct{
	*Queries
	DB *sqlx.DB
}

func LoadRepo(){
	appSettings, err := appsettings.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", "Error getting app settings", logger.LOG_LEVEL_ERROR)
		panic("")
	}
	
	fmt.Println("Connecting to database --- ",appSettings.DBMasterConnectionString);
	connection,err := sqlx.Connect("postgres",appSettings.DBMasterConnectionString)
	if err!=nil{
		log.Fatal("Cannot connect to postgres database database ",err)
	}
	REPO=&Repo{
		Queries: New(connection),
		DB: connection,
	}
}