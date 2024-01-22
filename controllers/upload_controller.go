package controllers

import (
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type UploadController struct{}

func (uploadController UploadController) SaveToUploadsTable(file string, table string, columnID int32) {

	// CREATE TABLE uploads(
	// 	id SERIAL PRIMARY KEY,
	// 	item_id INTEGER,
	// 	type VARCHAR(100),
	// 	path TEXT,
	// 	related_table VARCHAR(150),
	// 	meta JSON NULL
	// );

	_, err := gen.REPO.DB.NamedExec(`INSERT INTO uploads (item_id,type,path,related_table) values(:item_id,:type,:path,:related_table) on conflict (item_id, related_table) do update set path = EXCLUDED.path`,
		map[string]interface{}{
			"item_id":       columnID,
			"type":          "file",
			"path":          file,
			"related_table": table,
		})
		if err!=nil{
			logger.Log("UploadController",err.Error(),logger.LOG_LEVEL_ERROR)
		}

}
