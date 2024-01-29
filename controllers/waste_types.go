package controllers

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"
	"time"
	"ttnmwastemanagementsystem/gen"

	"github.com/gin-gonic/gin"
)

type WasteTypesController struct{}

type InsertWasteGroupParams struct {
	Name     string `json:"name"  binding:"required"`
	Category string `json:"category"`
	ParentID int32  `json:"parent_id"`
	FilePath string `json:"file_path" binding:"required"`
}

type UpdateWasteGroupParams struct {
	ID       int    `json:"id"  binding:"required"`
	Name     string `json:"name"  binding:"required"`
	Category string `json:"category"`
	ParentID int32  `json:"parent_id"`
	IsActive *bool  `json:"is_active"`
	FilePath string `json:"file_path"`
}

func (wasteGroupsController WasteTypesController) InsertWasteGroup(context *gin.Context) {
	var params InsertWasteGroupParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	wasteType, insertError := gen.REPO.InsertWasteType(context, gen.InsertWasteTypeParams{
		Name:     params.Name,
		ParentID: sql.NullInt32{Int32: params.ParentID, Valid: params.ParentID != 0},
	})

	if insertError != nil {
		fmt.Println(insertError.Error())
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to add Waste Type",
		})
		return
	}

	UploadController{}.SaveToUploadsTable(params.FilePath, "waste_types", wasteType.ID)
	// If you want to return the created company as part of the response
	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Successfully Created Waste Type",
		"waste_type": wasteType, // Include the company details in the response
	})

}

func (wasteGroupsController WasteTypesController) GetAllWasteTypes(context *gin.Context) {
	parentWasteTypeFilter := context.Query("p")
	parentWasteTypeFilter_, _ := strconv.ParseUint(parentWasteTypeFilter, 10, 32)

	type Result struct {
		ID        int32          `json:"id"`
		Name      string         `json:"name"`
		IsActive  bool           `json:"is_active"`
		ParentID  sql.NullInt32  `json:"parent_id"`
		CreatedAt time.Time      `json:"created_at"`
		FilePath  sql.NullString `json:"file_path"`
		Children  []Result       `json:"children"`
	}

	//results := []Result{}
	if parentWasteTypeFilter != "" {
		wasteTypes, err := gen.REPO.GetAllWasteTypes(context)
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"error":       false,
			"waste_types": wasteTypes,
		})
	} else {
		wasteTypes, err := gen.REPO.GetChildrenWasteTypes(context, sql.NullInt32{Int32: int32(parentWasteTypeFilter_), Valid: true})
		if err != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"error":       false,
			"waste_types": wasteTypes,
		})

	}

	// for _, v := range mainWasteTypes {
	// 	result := Result{}
	// 	result.ID = v.ID
	// 	result.Name = v.Name
	// 	result.IsActive = v.IsActive
	// 	result.ParentID = v.ParentID
	// 	result.CreatedAt = v.CreatedAt
	// 	result.FilePath = v.FilePath

	// 	children := []Result{}
	// 	fmt.Println(v.ParentID.Int32)
	// 	childrenWasteTypes, _ := gen.REPO.GetChildrenWasteTypes(context, sql.NullInt32{Int32: v.ID, Valid: true})
	// 	for _, x := range childrenWasteTypes {
	// 		child := Result{}
	// 		child.ID = x.ID
	// 		child.Name = x.Name
	// 		child.IsActive = x.IsActive
	// 		child.ParentID = x.ParentID
	// 		child.CreatedAt = x.CreatedAt
	// 		child.FilePath = x.FilePath

	// 		children = append(children, child)
	// 	}
	// 	result.Children = children
	// 	results = append(results, result)
	// }

	// context.JSON(http.StatusOK, gin.H{
	// 	"error":       false,
	// 	"waste_types": results,
	// })
}

func (wasteGroupsController WasteTypesController) GetUsersWasteGroups(context *gin.Context) {
	wasteGroups, err := gen.REPO.GetUsersWasteType(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"waste_types": wasteGroups,
	})
}

func (wasteGroupController WasteTypesController) GetOneWasteGroup(context *gin.Context) {
	id := context.Param("id")

	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	wasteGroup, err := gen.REPO.GetOneWasteType(context, int32(id_))

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":      false,
		"waste_type": wasteGroup,
	})
}

func (wasteGroupController WasteTypesController) UpdateWasteType(context *gin.Context) {
	var params UpdateWasteGroupParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update Waste Group
	updateError := gen.REPO.UpdateWasteType(context, gen.UpdateWasteTypeParams{
		Name:     params.Name,
		ParentID: sql.NullInt32{Int32: params.ParentID, Valid: params.ParentID != 0},
		ID:       int32(params.ID),
		IsActive: *params.IsActive,
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	UploadController{}.SaveToUploadsTable(params.FilePath, "waste_types", int32(params.ID))

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated Waste Group",
	})
}
