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

type FileController struct{}

func (fileController FileController) UploadFiles(context *gin.Context) {
	form, _ := context.MultipartForm()
	files := form.File["upload[]"]
	filePaths := []string{}
	for _, file := range files {
		logger.Log("UploadController", fmt.Sprint("Uploading file::", file.Filename), logger.LOG_LEVEL_INFO)
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
func (controller FileController) GetFile(context *gin.Context) {
	file, _ := context.Params.Get("file")
	basePath, _ := os.Getwd()

	completePath := basePath + "/uploads/" + file

	_, err := os.Stat(completePath)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "File not found",
		})
		return
	}
	context.File(completePath)
}

func (controller FileController) GetFlag(context *gin.Context) {
	file, _ := context.Params.Get("file")
	basePath, _ := os.Getwd()

	completePath := basePath + "/assets/flags/" + file

	_, err := os.Stat(completePath)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "File not found",
		})
		return
	}
	context.File(completePath)
}


func DeleteFile(filePath string) {

}
