package models

type Location struct {
	City    string  `db:"home_city" json:"city"`
	State   string  `db:"home_state" json:"state"`
	Zip     string  `db:"home_zip" json:"zip"`
	Country string  `db:"home_country" json:"country"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}
