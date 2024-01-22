package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/utils"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func (uploa UploadController) UploadFiles(context *gin.Context) {
	form, _ := context.MultipartForm()
	files := form.File["upload[]"]
	filePaths := []string{}

	fmt.Println(files)
	for _, file := range files {
		logger.Log("UploadController", fmt.Sprint("Uploading file::", file.Filename), logger.LOG_LEVEL_INFO)

		strippedFileName:= strings.ReplaceAll(file.Filename, " ", "")
		parsedFileName := fmt.Sprint(utils.NewRandStringLen(50), strippedFileName)
		var base = filepath.Base("/files/"+parsedFileName)

		logger.Log("UploadController",fmt.Sprint("Saving file to ::",base),logger.LOG_LEVEL_INFO)
		filePaths = append(filePaths, base)
		context.SaveUploadedFile(file, base)
	}
	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"file_paths": filePaths,
	})
}
