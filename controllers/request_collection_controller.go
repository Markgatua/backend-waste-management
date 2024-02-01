package controllers

import (
	"database/sql"
	_"encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	_ "ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/models"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-module/carbon"
	"github.com/guregu/null"
	_ "gopkg.in/guregu/null.v3"
)

type RequestCollectionController struct{}

type InsertNewCollectionRequestParams struct {
	ProducerID     int32 `json:"producer_id"  binding:"required"`
	RequestDate    time.Time `json:"request_date"  binding:"required"`
	Location       models.Location `json:"location" binding:"required"`
	ContactPerson  string `json:"contact_person"  binding:"required"`
}

type InsertNewNotificationRequestParams struct {
	UserID  int32 `json:"user_id"  binding:"required"`
	Subject string        `json:"subject"  binding:"required"`
	Message string        `json:"message"  binding:"required"`
}

type ConfirmCollectionRequestParams struct {
	Confirm *bool 			`json:"confirm"`
	ID        int32        `json:"id"`
}

type CancelCollectionRequestParams struct {
	Cancel *bool `json:"cancel"`
	ID        int32        `json:"id"`
}

type UpdateCollectionRequestParams struct {
	PickupDate time.Time `json:"pickup_date"   binding:"required"`
	Status     *bool `json:"status"`
	ID         int32        `json:"id"`
}

func (requestCollectionController RequestCollectionController) InsertNewCollectionRequestParams(context *gin.Context) {
	var params InsertNewCollectionRequestParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	championIDNullable := int32(params.ProducerID)

	ChampionCollector, err := gen.REPO.GetTheCollectorForAChampion(context, championIDNullable)

	championCID := ChampionCollector.CollectorID

	// var params2 InsertNewNotificationRequestParams

	 insertError := gen.REPO.InsertNewCollectionRequest(context, gen.InsertNewCollectionRequestParams{
		ProducerID: params.ProducerID,
		CollectorID: championCID,
		RequestDate: params.RequestDate,
		Location: null.StringFrom(params.Location.Location).NullString,
		AdministrativeLevel1Location: null.StringFrom(params.Location.AdministrativeAreaLevel1).NullString,
		Lat: null.FloatFrom(params.Location.LatLng.Lat).NullFloat64,
		Lng: null.FloatFrom(params.Location.LatLng.Lng).NullFloat64,
		FirstContactPerson: params.ContactPerson,
	})

	fmt.Println(insertError);

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Send Your Collection Request",
		})
		return
	}

	notificationUsers, _ := gen.REPO.GetCompanyUsers(context, sql.NullInt32{Int32: championCID,Valid: true})

	producerData, _ := gen.REPO.GetCompany(context, params.ProducerID)
	var subject = producerData.Name +" "+ producerData.Location.String;

	fmt.Println(subject)


	// Parse the date string using Carbon
	parsedTime,err := time.Parse("2006-01-02 15:04:05 -0700 MST", params.RequestDate.String())
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// Format the date in the desired format
	formattedDate := parsedTime.Format("Jan 2, 2006 03:04 pm")
	

	fmt.Println(formattedDate);

	for _, value := range notificationUsers {

		insertErroh := gen.REPO.InsertNewNotificationRequest(context, gen.InsertNewNotificationRequestParams{
			UserID: value.ID,
			Subject: subject,
			Message: "You have a new Collection Request",
		})

		phoneNumber :=value.CallingCode.String + value.Phone.String
		phone := phoneNumber
		sms := helpers.SMS{}
		err = sms.SendSMS([]string{phone}, fmt.Sprint("You Have a new Collection Request from " + subject + ".\n\n Collection Date: " + formattedDate ))

		fmt.Println(insertErroh);
    }

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully sent Your Collection Request",
	})

}

func (requestCollectionController RequestCollectionController) ConfirmCollectionRequest(context *gin.Context) {
	var params ConfirmCollectionRequestParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	insertError := gen.REPO.ConfirmCollectionRequest(context, gen.ConfirmCollectionRequestParams{
		ID:        params.ID,
		Confirmed: sql.NullBool{Bool: *params.Confirm,Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Confirm Collection Request",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Confirmed the Collection Request",
	})
}

func (requestCollectionController RequestCollectionController) CancelCollectionRequest(context *gin.Context) {
	var params CancelCollectionRequestParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	insertError := gen.REPO.CancelCollectionRequest(context, gen.CancelCollectionRequestParams{
		ID:        params.ID,
		Cancelled: sql.NullBool{Bool: *params.Cancel, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Cancel Collection Request",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Cancelled the Collection Request",
	})
}

func (requestCollectionController RequestCollectionController) UpdateCollectionRequest(context *gin.Context) {
	
	var params UpdateCollectionRequestParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	insertError := gen.REPO.UpdateCollectionRequest(context, gen.UpdateCollectionRequestParams{
		ID:        params.ID,
		Status: sql.NullBool{Bool: *params.Status, Valid: true},
		PickupDate: sql.NullTime{Time: params.PickupDate, Valid: true},
	})

	if insertError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
	})

}

func (requestCollectionController RequestCollectionController) CollectionWeightTotals(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	
	count,insertError := gen.REPO.CollectionWeightTotals(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}


func (requestCollectionController RequestCollectionController) GetLatestCollection(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);


	idp, insertErrorh := gen.REPO.GetProducerLatestCollectionId(context,id32);

	fmt.Println(insertErrorh)

	
	count,insertError := gen.REPO.GetLatestCollection(context,idp.ID)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}


func (requestCollectionController RequestCollectionController) GetProducerLatestCollectionId(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	count,insertError := gen.REPO.GetLatestCollection(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}



func (requestCollectionController RequestCollectionController) GetWasteItemsProducerData(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	count,insertError := gen.REPO.GetWasteItemsProducerData(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}

func (requestCollectionController RequestCollectionController) GetCollectionStats(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	count,insertError := gen.REPO.GetCollectionStats(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}

func (requestCollectionController RequestCollectionController) GetAllProducerCompletedCollectionRequests(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	count,insertError := gen.REPO.GetAllProducerCompletedCollectionRequests(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}

func (requestCollectionController RequestCollectionController) GetAllProducerPendingCollectionRequests(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	count,insertError := gen.REPO.GetAllProducerPendingCollectionRequests(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Update Collection Request",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	fmt.Println(count)


	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Updated the Collection Request",
		"count": count,
	})
}



func (requestCollectionController RequestCollectionController) GetTheCollectorForAChampion(context *gin.Context) {
	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************")
	fmt.Println(id32)

	data, insertError := gen.REPO.GetTheCollectorForAChampion(context, id32)

	if insertError != nil {
		logger.Log("RequestCollectionController", insertError.Error(), logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Fetch my next collection Date",
		})
		return
	} else {
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	type ChampionAggregatorAssignment struct {
		PickupDay   sql.NullString `json:"pickup_day"`
		PickupTime  sql.NullString `json:"pickup_time"`
		NextDate string `json:"next_date"`
		Collector string `json:"collector"`
	}

	

	fmt.Println(data)

	pickupDay := data.PickupDay.String
	valid := data.PickupDay.Valid

	if !valid {
		fmt.Println("Error accessing pickup_day.String field in data.")
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Internal Server Error",
		})
		return
	}


	// Get the current day
	currentDay := time.Now().Weekday().String()

	if currentDay == pickupDay {
		data2:= ChampionAggregatorAssignment{}
		data2.PickupTime=data.PickupTime
		data2.PickupDay=data.PickupDay
		data2.NextDate="Today"
		data2.Collector=data.CollectorName.String

		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Collection day is today!",
			"data":    data2,
		})
	} else {
		// Calculate days until next Monday
		daysUntilNextMonday := (int(time.Monday) - int(time.Now().Weekday()) + 7) % 7

		// Calculate the next Monday date
		nextMonday := time.Now().AddDate(0, 0, daysUntilNextMonday)

		data2:= ChampionAggregatorAssignment{}
		data2.PickupTime=data.PickupTime
		data2.PickupDay=data.PickupDay
		data2.NextDate=nextMonday.Format("02-01-2006")
		data2.Collector=data.CollectorName.String


		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Next collection day is on " + nextMonday.Format("02-01-2006"),
			"data":    data2,
		})
	}
}

func (requestCollectionController RequestCollectionController) GetAggregatorNewRequests(context *gin.Context){

	id, _ := context.Params.Get("id")
	var id32 int32
	fmt.Sscan(id, &id32)

	fmt.Println("********************");
	fmt.Println(id32);

	
	data,insertError := gen.REPO.GetAggregatorNewRequests(context,id32)

	if insertError != nil {
		logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Failed to Fetch Aggregator New Requests",
		})
		return
	}else{
		//logger.Log("RequestCollectionController",insertError.Error(),logger.LOG_LEVEL_ERROR)
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Successfully Fetched Aggregators new Requests",
		"data": data,
	})
}