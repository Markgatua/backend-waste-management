package models

import "time"

type TtnmOrganizationModel struct {
	ID        int64     `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	AboutUs     string    `db:"about_us" json:"about_us"`
	TagLine     string    `db:"tag_line" json:"tag_line"`
	LogoPath     string    `db:"logo_path" json:"logo_path"`
	WebsiteUrl     string    `db:"website_url" json:"website_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}