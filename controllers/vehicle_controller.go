package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
)

type VehicleController struct{}

type CreateVehicleParams struct {
	VehicleTypeID    int32  `json:"vehicle_type_id" binding:"required"`
	RegNo            string `json:"reg_no" binding:"required"`
	AssignedDriverID *int32 `json:"assigned_driver_id"`
	IsActive         *bool  `json:"is_active"  binding:"required"`
}

type UpdateVehicleParams struct {
	ID               int    `json:"id"  binding:"required"`
	VehicleTypeID    int32  `json:"vehicle_type_id" binding:"required"`
	RegNo            string `json:"reg_no" binding:"required"`
	AssignedDriverID *int32 `json:"assigned_driver_id"`
	IsActive         *bool  `json:"is_active"  binding:"required"`
}

type UpdateVehicleStatusParams struct {
	ID       int   `json:"id"  binding:"required"`
	IsActive *bool `json:"status"  binding:"required"`
}

func (controller VehicleController) InsertVehicle(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)

	var params CreateVehicleParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	count, err := gen.REPO.GetDuplicateVehicle(context, params.RegNo)
	if err != nil && err != sql.ErrNoRows {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if count != 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Vehicle exists",
		})
		return
	}
	if params.AssignedDriverID != nil {
		vehicle, insertError := gen.REPO.AddVehicle(context, gen.AddVehicleParams{
			CompanyID:        int32(auth.UserCompanyId.Int64),
			AssignedDriverID: sql.NullInt32{Int32: *params.AssignedDriverID, Valid: params.AssignedDriverID != nil},
			VehicleTypeID:    params.VehicleTypeID,
			RegNo:            params.RegNo,
			IsActive:         *params.IsActive,
		})
		if insertError != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Failed to add vehicle",
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Successfully Added vehicle",
			"vehicle": vehicle,
		})
	} else {
		vehicle, insertError := gen.REPO.AddVehicle(context, gen.AddVehicleParams{
			CompanyID:     int32(auth.UserCompanyId.Int64),
			VehicleTypeID: params.VehicleTypeID,
			RegNo:         params.RegNo,
			IsActive:      *params.IsActive,
		})
		if insertError != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Failed to add vehicle",
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Successfully Added vehicle",
			"vehicle": vehicle,
		})
	}

}

func (controller VehicleController) GetAllVehicleTypes(context *gin.Context) {
	items, err := gen.REPO.GetVehicleTypes(context)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": items,
	})
}
func (controller VehicleController) GetAllVehicles(context *gin.Context) {
	auth, _ := helpers.Functions{}.CurrentUserFromToken(context)
	items, err := gen.REPO.GetAllVehicles(context, int32(auth.UserCompanyId.Int64))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"content": items,
	})
}

func (c VehicleController) DeleteVehicle(context *gin.Context) {
	id := context.Param("id")
	id_, _ := strconv.ParseUint(id, 10, 32)
	println("------------------------------", id_)
	err := gen.REPO.DeleteVehicle(context, int32(id_))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Vehicle successfully deleted",
	})
}

func (controller VehicleController) UpdateVehicleStatus(context *gin.Context) {
	var params UpdateVehicleStatusParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update company status
	updateError := gen.REPO.UpdateVehicleStatus(context, gen.UpdateVehicleStatusParams{
		IsActive: *params.IsActive,
		ID:       int32(params.ID),
	})
	if updateError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": updateError.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated vehicle",
		"status":  params.IsActive, // Use the variable for the parsed status
	})
}

func (controller VehicleController) UpdateVehicle(context *gin.Context) {
	var params UpdateVehicleParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	count, err := gen.REPO.GetDuplicateVehiclesWithoutID(context, gen.GetDuplicateVehiclesWithoutIDParams{
		ID:    int32(params.ID),
		RegNo: strings.ToLower(params.RegNo),
	})
	if count > 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Another vehicle has the same reg no",
		})
		return
	}

	if params.AssignedDriverID != nil {
		updateError := gen.REPO.UpdateVehicle(context, gen.UpdateVehicleParams{
			IsActive:         *params.IsActive,
			RegNo:            params.RegNo,
			VehicleTypeID:    params.VehicleTypeID,
			AssignedDriverID: sql.NullInt32{Int32: *params.AssignedDriverID, Valid: params.AssignedDriverID != nil},
			ID:               int32(params.ID),
		})
		if updateError != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": updateError.Error(),
			})
			return
		}
	} else {
		updateError := gen.REPO.UpdateVehicle(context, gen.UpdateVehicleParams{
			IsActive:         *params.IsActive,
			RegNo:            params.RegNo,
			VehicleTypeID:    params.VehicleTypeID,
			ID:               int32(params.ID),
		})
		if updateError != nil {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": updateError.Error(),
			})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully updated vehicle",
	})
}
