package geopoint

import "time"

type GeoPoint struct {
	Id         int       `db:"id" json:"id" example:"1"`
	Title      string    `db:"title" json:"title" example:"Forêt à l'aube"`
	UserId     int       `db:"user_id" json:"user_id" example:"1"`
	Location   Location  `db:"location" json:"location"`
	CreatedOn  time.Time `db:"created_on" json:"created_on" example:"2022-05-26T11:17:35.079344Z"`
	Amplitudes []int     `db:"amplitudes" json:"amplitudes" example:"0,1,2,3,45,3,2,1"`
	Picture    string    `db:"picture" json:"picture" example:"https://example.com/picture-1.jpg"`
	Sound      string    `db:"sound" json:"sound" example:"https://example.com/sound-2.mp3"`
}

type Location struct {
	Todo string
}

type AddGeoPoint struct {
	Title      string    `json:"title" example:"Forêt à l'aube"`
	UserId     int       `json:"user_id" example:"1"`
	Location   []float64 `json:"location" example:"38.652608,-120.357448"`
	Date       time.Time `json:"date" example:"2022-05-26T11:17:35.079344Z"`
	Amplitudes []int     `json:"amplitudes" example:"0,1,2,3,45,3,2,1"`
}
