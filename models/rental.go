package models

import "time"

type RentalPrice struct {
	Day int `db:"price_per_day" json:"day"`
}

type Rental struct {
	Id              int     `json:"id"`
	UserID          int     `db:"user_id" json:"-"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Make            string  `db:"vehicle_make" json:"make"`
	Model           string  `db:"vehicle_model" json:"model"`
	Year            int     `db:"vehicle_year" json:"year"`
	Length          float32 `db:"vehicle_length" json:"length"`
	Sleeps          int     `json:"sleeps"`
	PrimaryImageURL string  `db:"primary_image_url" json:"primary_image_url"`
	RentalPrice     `json:"rental_price"`
	Location        `json:"location"`
	User            `json:"user"`
	Created         time.Time `json:"-"`
	Updated         time.Time `json:"-"`
}
