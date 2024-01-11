package models

import "time"

type EmailVerificationToken struct {
	ID        int64     `db:"id" json:"id"`
	Token     string    `db:"token" json:"token"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}