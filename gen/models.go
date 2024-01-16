// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package gen

import (
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

type ChampionAggregatorAssignment struct {
	ID          int32         `json:"id"`
	ChampionID  sql.NullInt32 `json:"champion_id"`
	CollectorID sql.NullInt32 `json:"collector_id"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Company struct {
	ID             int32          `json:"id"`
	Name           string         `json:"name"`
	CompanyType    int32          `json:"company_type"`
	OrganizationID int32          `json:"organization_id"`
	Region         sql.NullString `json:"region"`
	Location       sql.NullString `json:"location"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
}

type Country struct {
	ID               int32          `json:"id"`
	Name             string         `json:"name"`
	CurrencyCode     sql.NullString `json:"currency_code"`
	Capital          sql.NullString `json:"capital"`
	Citizenship      string         `json:"citizenship"`
	CountryCode      string         `json:"country_code"`
	Currency         sql.NullString `json:"currency"`
	CurrencySubUnit  sql.NullString `json:"currency_sub_unit"`
	CurrencySymbol   sql.NullString `json:"currency_symbol"`
	CurrencyDecimals sql.NullInt16  `json:"currency_decimals"`
	FullName         sql.NullString `json:"full_name"`
	Iso31662         string         `json:"iso_3166_2"`
	Iso31663         string         `json:"iso_3166_3"`
	RegionCode       string         `json:"region_code"`
	SubRegionCode    string         `json:"sub_region_code"`
	Eea              sql.NullInt16  `json:"eea"`
	CallingCode      sql.NullString `json:"calling_code"`
	Flag             sql.NullString `json:"flag"`
}

type EmailVerificationToken struct {
	ID        int32     `json:"id"`
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Organization struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	CountryID int32  `json:"country_id"`
}

type PaymentMethod struct {
	ID            int32     `json:"id"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

type Permission struct {
	ID        int32          `json:"id"`
	Name      string         `json:"name"`
	Action    string         `json:"action"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	Module    string         `json:"module"`
	Submodule sql.NullString `json:"submodule"`
}

type PhoneVerificationToken struct {
	ID          int32     `json:"id"`
	Token       string    `json:"token"`
	CallingCode string    `json:"calling_code"`
	Phone       string    `json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
}

type Role struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	GuardName   string         `json:"guard_name"`
	CreatedAt   sql.NullTime   `json:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at"`
	Description sql.NullString `json:"description"`
	IsActive    bool           `json:"is_active"`
}

type RoleHasPermission struct {
	PermissionID int32 `json:"permission_id"`
	RoleID       int32 `json:"role_id"`
}

type TtnmOrganization struct {
	ID                     int32     `json:"id"`
	OrganizationID         string    `json:"organization_id"`
	Name                   string    `json:"name"`
	TagLine                string    `json:"tag_line"`
	AboutUs                string    `json:"about_us"`
	LogoPath               string    `json:"logo_path"`
	AppAppstoreLink        string    `json:"app_appstore_link"`
	AppGooglePlaystoreLink string    `json:"app_google_playstore_link"`
	WebsiteUrl             string    `json:"website_url"`
	City                   string    `json:"city"`
	State                  string    `json:"state"`
	Zip                    string    `json:"zip"`
	Country                string    `json:"country"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type Upload struct {
	ID           int32                 `json:"id"`
	ItemID       sql.NullInt32         `json:"item_id"`
	Type         sql.NullString        `json:"type"`
	Path         sql.NullString        `json:"path"`
	RelatedTable sql.NullString        `json:"related_table"`
	Meta         pqtype.NullRawMessage `json:"meta"`
}

type User struct {
	ID                 int32          `json:"id"`
	FirstName          sql.NullString `json:"first_name"`
	LastName           sql.NullString `json:"last_name"`
	Provider           sql.NullString `json:"provider"`
	RoleID             sql.NullInt32  `json:"role_id"`
	UserCompanyID      sql.NullInt32  `json:"user_company_id"`
	Email              sql.NullString `json:"email"`
	Password           sql.NullString `json:"password"`
	AvatarUrl          sql.NullString `json:"avatar_url"`
	UserType           sql.NullInt16  `json:"user_type"`
	IsActive           sql.NullBool   `json:"is_active"`
	CallingCode        sql.NullString `json:"calling_code"`
	Phone              sql.NullString `json:"phone"`
	PhoneConfirmedAt   sql.NullTime   `json:"phone_confirmed_at"`
	ConfirmedAt        sql.NullTime   `json:"confirmed_at"`
	ConfirmationToken  sql.NullString `json:"confirmation_token"`
	ConfirmationSentAt sql.NullTime   `json:"confirmation_sent_at"`
	RecoveryToken      sql.NullString `json:"recovery_token"`
	RecoverySentAt     sql.NullTime   `json:"recovery_sent_at"`
	LastLogin          sql.NullTime   `json:"last_login"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type WasteBuyer struct {
	ID      int32                 `json:"id"`
	BuyerID sql.NullInt32         `json:"buyer_id"`
	Rates   pqtype.NullRawMessage `json:"rates"`
}

type WasteCollection struct {
	ID          int32                 `json:"id"`
	Date        time.Time             `json:"date"`
	ChampionID  sql.NullInt32         `json:"champion_id"`
	CollectorID sql.NullInt32         `json:"collector_id"`
	Waste       pqtype.NullRawMessage `json:"waste"`
	IsCollected sql.NullBool          `json:"is_collected"`
	CreatedAt   time.Time             `json:"created_at"`
}

type WasteForSale struct {
	ID     int32                 `json:"id"`
	Seller sql.NullInt32         `json:"seller"`
	Waste  pqtype.NullRawMessage `json:"waste"`
}

type WasteGroup struct {
	ID        int32        `json:"id"`
	Name      string       `json:"name"`
	Category  string       `json:"category"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type WasteTransaction struct {
	ID                int32                 `json:"id"`
	Date              time.Time             `json:"date"`
	BuyerID           sql.NullInt32         `json:"buyer_id"`
	SellerID          sql.NullInt32         `json:"seller_id"`
	WasteProducts     pqtype.NullRawMessage `json:"waste_products"`
	TotalAmount       string                `json:"total_amount"`
	PaymentMethodID   sql.NullInt32         `json:"payment_method_id"`
	MerchantRequestID sql.NullString        `json:"merchant_request_id"`
	CheckoutRequestID sql.NullString        `json:"checkout_request_id"`
	MpesaResultCode   sql.NullString        `json:"mpesa_result_code"`
	MpesaResultDesc   sql.NullString        `json:"mpesa_result_desc"`
	MpesaReceiptCode  sql.NullString        `json:"mpesa_receipt_code"`
	TimePaid          sql.NullString        `json:"time_paid"`
	IsPaid            sql.NullBool          `json:"is_paid"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}
