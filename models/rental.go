package models

type RentalPrice struct {
	Day int `json:"day"`
}

type Rental struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Type            string      `json:"type"`
	Make            string      `json:"make"`
	Model           string      `json:"model"`
	Year            int         `json:"year"`
	Length          float32     `json:"length"`
	Sleeps          int         `json:"sleeps"`
	PrimaryImageURL string      `json:"primary_image_url"`
	Price           RentalPrice `json:"price"`
	Location        Location    `json:"location"`
}
