package seeder

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"
)

func Run() {
	appSettings, err := utils.GetAppSettings()
	if err != nil {
		logger.Log("SEEDER", "Error getting app settings", logger.LOG_LEVEL_ERROR)
		return
	}

	conn, err := pgx.Connect(context.Background(), appSettings.DBMasterConnectionString)
	defer conn.Close(context.Background())
	if err != nil {
		logger.Log("SEEDER", fmt.Sprint("Unable to connect to database: %v", err), logger.LOG_LEVEL_ERROR)
	} else {
		queries := New(conn)
	}

}

