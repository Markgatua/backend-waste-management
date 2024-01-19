package models

import "gopkg.in/guregu/null.v3"

type User struct {
	ID                 null.Int    `db:"id" json:"id"`
	FirstName          null.String `db:"first_name" json:"first_name"`
	LastName           null.String `db:"last_name" json:"last_name"`
	Provider           null.String `db:"provider" json:"provider"`
	Email              null.String `db:"email" json:"email"`
	RoleId             null.Int    `db:"role_id" json:"role_id"`
	UserCompanyId      null.Int    `db:"user_company_id"  json:"user_company_id"`
	Password           null.String `db:"password" json:"-"`
	ISttnm_user        bool        `db:"is_ttnm_user" json:"-"`
	IsActive           null.Bool   `db:"is_active" json:"is_active"`
	CallingCode        null.String `db:"calling_code" json:"calling_code"`
	AvatarUrl          null.String `db:"avatar_url" json:"avatar_url"`
	Phone              null.String `db:"phone" json:"phone"`
	PhoneConfirmedAt   null.Time   `db:"phone_confirmed_at" json:"phone_confirmed_at"`
	ConfirmedAt        null.Time   `db:"confirmed_at" json:"confirmed_at"`
	ConfirmationToken  null.String `db:"confirmation_token" json:"-"`
	ConfirmationSentAt null.Time   `db:"confirmation_sent_at" json:"confirmation_sent_at"`
	RecoveryToken      null.String `db:"recovery_token" json:"-"`
	RecoverySentAt     null.Time   `db:"recovery_sent_at" json:"-"`
	CreatedAt          null.Time   `db:"created_at" json:"created_at"`
	UpdatedAt          null.Time   `db:"updated_at" json:"-"`
	UserType           null.Int    `db:"user_type" json:"-"`
	LastLogin          null.Time   `db:"last_login" json:"last_login"`
}
