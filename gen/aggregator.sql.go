// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: aggregator.sql

package gen

import (
	"context"
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

const createBuyer = `-- name: CreateBuyer :one
insert into
    buyers (
        company_id, company, first_name, last_name, calling_code, phone, administrative_level_1_location, location, is_active, lat, lng, created_at, updated_at,region
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14
    ) returning id, company_id, company, first_name, last_name, is_active, region, calling_code, location, administrative_level_1_location, lat, lng, phone, created_at, updated_at
`

type CreateBuyerParams struct {
	CompanyID                    int32           `json:"company_id"`
	Company                      sql.NullString  `json:"company"`
	FirstName                    string          `json:"first_name"`
	LastName                     string          `json:"last_name"`
	CallingCode                  sql.NullString  `json:"calling_code"`
	Phone                        sql.NullString  `json:"phone"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	Location                     sql.NullString  `json:"location"`
	IsActive                     bool            `json:"is_active"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	CreatedAt                    time.Time       `json:"created_at"`
	UpdatedAt                    time.Time       `json:"updated_at"`
	Region                       sql.NullString  `json:"region"`
}

func (q *Queries) CreateBuyer(ctx context.Context, arg CreateBuyerParams) (Buyer, error) {
	row := q.db.QueryRowContext(ctx, createBuyer,
		arg.CompanyID,
		arg.Company,
		arg.FirstName,
		arg.LastName,
		arg.CallingCode,
		arg.Phone,
		arg.AdministrativeLevel1Location,
		arg.Location,
		arg.IsActive,
		arg.Lat,
		arg.Lng,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Region,
	)
	var i Buyer
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.Company,
		&i.FirstName,
		&i.LastName,
		&i.IsActive,
		&i.Region,
		&i.CallingCode,
		&i.Location,
		&i.AdministrativeLevel1Location,
		&i.Lat,
		&i.Lng,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createSale = `-- name: CreateSale :one
insert into sales(ref,company_id,buyer_id,total_amount_of_waste,total_amount,dump) VALUES ($1,$2,$3,$4,$5,$6) returning id, ref, company_id, buyer_id, total_amount_of_waste, total_amount, date, dump
`

type CreateSaleParams struct {
	Ref                string                `json:"ref"`
	CompanyID          int32                 `json:"company_id"`
	BuyerID            int32                 `json:"buyer_id"`
	TotalAmountOfWaste sql.NullString        `json:"total_amount_of_waste"`
	TotalAmount        sql.NullString        `json:"total_amount"`
	Dump               pqtype.NullRawMessage `json:"dump"`
}

func (q *Queries) CreateSale(ctx context.Context, arg CreateSaleParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, createSale,
		arg.Ref,
		arg.CompanyID,
		arg.BuyerID,
		arg.TotalAmountOfWaste,
		arg.TotalAmount,
		arg.Dump,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.Ref,
		&i.CompanyID,
		&i.BuyerID,
		&i.TotalAmountOfWaste,
		&i.TotalAmount,
		&i.Date,
		&i.Dump,
	)
	return i, err
}

const createSaleItem = `-- name: CreateSaleItem :one
insert into sale_items(company_id,sale_id,waste_type_id,amount_of_waste,cost_per_kg,total_amount) VALUES($1,$2,$3,$4,$5,$6) returning id, company_id, sale_id, waste_type_id, amount_of_waste, cost_per_kg, total_amount
`

type CreateSaleItemParams struct {
	CompanyID     int32          `json:"company_id"`
	SaleID        int32          `json:"sale_id"`
	WasteTypeID   int32          `json:"waste_type_id"`
	AmountOfWaste sql.NullString `json:"amount_of_waste"`
	CostPerKg     sql.NullString `json:"cost_per_kg"`
	TotalAmount   string         `json:"total_amount"`
}

func (q *Queries) CreateSaleItem(ctx context.Context, arg CreateSaleItemParams) (SaleItem, error) {
	row := q.db.QueryRowContext(ctx, createSaleItem,
		arg.CompanyID,
		arg.SaleID,
		arg.WasteTypeID,
		arg.AmountOfWaste,
		arg.CostPerKg,
		arg.TotalAmount,
	)
	var i SaleItem
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.SaleID,
		&i.WasteTypeID,
		&i.AmountOfWaste,
		&i.CostPerKg,
		&i.TotalAmount,
	)
	return i, err
}

const deleteBuyer = `-- name: DeleteBuyer :exec
delete from buyers where id = $1
`

func (q *Queries) DeleteBuyer(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteBuyer, id)
	return err
}

const makeCashPayment = `-- name: MakeCashPayment :one
insert into sale_transactions(ref,sale_id,company_id,payment_method,transaction_date) VALUES($1,$2,$3,"CASH",$4) returning ref, id, sale_id, company_id, payment_method, checkout_request_id, merchant_request_id, card_mask, msisdn_idnum, transaction_date, receipt_no, amount, mpesa_result_code, mpesa_result_desc, ipay_status, created_at, updated_at
`

type MakeCashPaymentParams struct {
	Ref             sql.NullString `json:"ref"`
	SaleID          int32          `json:"sale_id"`
	CompanyID       int32          `json:"company_id"`
	TransactionDate sql.NullTime   `json:"transaction_date"`
}

func (q *Queries) MakeCashPayment(ctx context.Context, arg MakeCashPaymentParams) (SaleTransaction, error) {
	row := q.db.QueryRowContext(ctx, makeCashPayment,
		arg.Ref,
		arg.SaleID,
		arg.CompanyID,
		arg.TransactionDate,
	)
	var i SaleTransaction
	err := row.Scan(
		&i.Ref,
		&i.ID,
		&i.SaleID,
		&i.CompanyID,
		&i.PaymentMethod,
		&i.CheckoutRequestID,
		&i.MerchantRequestID,
		&i.CardMask,
		&i.MsisdnIdnum,
		&i.TransactionDate,
		&i.ReceiptNo,
		&i.Amount,
		&i.MpesaResultCode,
		&i.MpesaResultDesc,
		&i.IpayStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBuyer = `-- name: UpdateBuyer :exec
update buyers set company_id=$1,company=$2,first_name=$3,last_name=$4,calling_code=$5,phone=$6,administrative_level_1_location=$7,location=$8,is_active=$9,lat=$10,lng=$11 ,region=$12 where id=$13
`

type UpdateBuyerParams struct {
	CompanyID                    int32           `json:"company_id"`
	Company                      sql.NullString  `json:"company"`
	FirstName                    string          `json:"first_name"`
	LastName                     string          `json:"last_name"`
	CallingCode                  sql.NullString  `json:"calling_code"`
	Phone                        sql.NullString  `json:"phone"`
	AdministrativeLevel1Location sql.NullString  `json:"administrative_level_1_location"`
	Location                     sql.NullString  `json:"location"`
	IsActive                     bool            `json:"is_active"`
	Lat                          sql.NullFloat64 `json:"lat"`
	Lng                          sql.NullFloat64 `json:"lng"`
	Region                       sql.NullString  `json:"region"`
	ID                           int32           `json:"id"`
}

func (q *Queries) UpdateBuyer(ctx context.Context, arg UpdateBuyerParams) error {
	_, err := q.db.ExecContext(ctx, updateBuyer,
		arg.CompanyID,
		arg.Company,
		arg.FirstName,
		arg.LastName,
		arg.CallingCode,
		arg.Phone,
		arg.AdministrativeLevel1Location,
		arg.Location,
		arg.IsActive,
		arg.Lat,
		arg.Lng,
		arg.Region,
		arg.ID,
	)
	return err
}
