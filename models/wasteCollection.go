// models.go

package models

import "database/sql"

// WasteCollectionReport represents the model for waste collection reports
type WasteCollectionReport struct {
	RequestID           int             `json:"request_id"`
	ProducerID          int             `json:"producer_id"`
	CollectorID         int             `json:"collector_id"`
	RequestDate         string          `json:"request_date"`
	PickupTimeStampID   int             `json:"pickup_time_stamp_id"`
	Location            string          `json:"location"`
	AdminLevel1Location string          `json:"administrative_level_1_location"`
	Lat                 float64         `json:"lat"`
	Lng                 float64         `json:"lng"`
	PickupDate          string          `json:"pickup_date"`
	Status              int             `json:"status"`
	FirstContactPerson  string          `json:"first_contact_person"`
	SecondContactPerson string          `json:"second_contact_person"`
	RequestCreatedAt    string          `json:"request_created_at"`
	WasteItemID         sql.NullInt64   `json:"waste_item_id"`
	WasteTypeID         sql.NullInt64   `json:"waste_type_id"`
	WasteCollectorID    sql.NullInt64   `json:"waste_collector_id"`
	Weight              sql.NullFloat64 `json:"weight"`
	WasteItemCreatedAt  sql.NullString  `json:"waste_item_created_at"`
}
