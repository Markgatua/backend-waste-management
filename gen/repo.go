package gen

import (
	"database/sql"
	"fmt"
	"log"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"

	// _"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var REPO *Repo

type Repo struct{
	*Queries
	DB *sql.DB
}

type Cat struct{
	Legs int
	FurColor string
}

func LoadRepo(){
	appSettings, err := utils.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", "Error getting app settings", logger.LOG_LEVEL_ERROR)
		panic("")
	}
	
	fmt.Println("Connecting to database --- ",appSettings.DBMasterConnectionString);
	connection,err := sql.Open("postgres",appSettings.DBMasterConnectionString)
	if err!=nil{
		log.Fatal("Cannot connect to postgres database database ",err)
	}
	REPO=&Repo{
		Queries: New(connection),
		DB: connection,
	}
}