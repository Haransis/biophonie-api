package geopoint

import (
	"mime/multipart"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/lib/pq"
)

type GeoPoint struct {
	Id         int           `db:"id" json:"id" example:"1"`
	Title      string        `db:"title" json:"title" example:"Forêt à l'aube"`
	UserId     int           `db:"user_id" json:"userId" example:"1"`
	Location   postgis.Point `db:"location" json:"location"`
	CreatedOn  time.Time     `db:"created_on" json:"createdOn" example:"2022-05-26T11:17:35.079344Z"`
	Amplitudes pq.Int64Array `db:"amplitudes" json:"amplitudes" swaggertype:"array,number" example:"0,1,2,3,45,3,2,1"`
	Picture    string        `db:"picture" json:"picture" example:"https://example.com/picture-1.jpg"`
	Sound      string        `db:"sound" json:"sound" example:"https://example.com/sound-2.wav"`
}

type AddGeoPoint struct {
	Title      string    `json:"title" example:"Forêt à l'aube" binding:"required"`
	UserId     int       `json:"userId" example:"1" binding:"required"`
	Location   []float64 `json:"location" example:"38.652608,-120.357448" binding:"required"`
	Date       time.Time `json:"date" example:"2022-05-26T11:17:35.079344Z" binding:"required"`
	Amplitudes []int64   `json:"amplitudes" example:"0,1,2,3,45,3,2,1" binding:"required"`
}

type BindGeopoint struct {
	Geopoint *multipart.FileHeader `form:"geopoint" binding:"required"`
	Sound    *multipart.FileHeader `form:"sound" binding:"required"`
	Picture  *multipart.FileHeader `form:"picture" binding:"required"`
}
