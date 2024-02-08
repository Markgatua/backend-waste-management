package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	_ "strings"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"

	"github.com/gin-gonic/gin"
)

type RoutePlanningController struct{}

func (controller RoutePlanningController) GetRoutes(context *gin.Context) {
	logger.Log("RoutePlanningController", fmt.Sprint("Running route planner"), logger.LOG_LEVEL_INFO)

	form, _ := context.MultipartForm()
	files := form.File["upload[]"]
	filePaths := []string{}
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		// Generate random file name for the new uploaded file so it doesn't override the old file with same name
		newFileName := utils.NewRandStringLen(50) + extension

		logger.Log("UploadController", fmt.Sprint("Saving file to ::", newFileName), logger.LOG_LEVEL_INFO)
		path, _ := os.Getwd()
		err := context.SaveUploadedFile(file, path+"/uploads/"+newFileName)
		if err != nil {
			logger.Log("UploadController", err.Error(), logger.LOG_LEVEL_ERROR)
		} else {
			filePaths = append(filePaths, newFileName)
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"file_paths": filePaths,
	})
}
