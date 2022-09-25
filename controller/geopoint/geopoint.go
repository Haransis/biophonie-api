package geopoint

import (
	"mime/multipart"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/lib/pq"
)

type GeoPoint struct {
	Id         int             `db:"id" json:"id" example:"1"`
	Title      string          `db:"title" json:"title" example:"Forêt à l'aube"`
	UserId     int             `db:"user_id" json:"userId" example:"1"`
	Latitude   float64         `json:"latitude"`
	Longitude  float64         `json:"longitude"`
	CreatedOn  time.Time       `db:"created_on" json:"createdOn" example:"2022-05-26T11:17:35.079344Z"`
	Amplitudes pq.Float64Array `db:"amplitudes" json:"amplitudes" swaggertype:"array,number" example:"0,1,2,3,45,3,2,1"`
	Picture    string          `db:"picture" json:"picture" example:"https://example.com/picture-1.jpg"`
	Sound      string          `db:"sound" json:"sound" example:"https://example.com/sound-2.wav"`
	Available  bool            `db:"available" json:"available" example:"true"`
}

type DbGeoPoint struct {
	*GeoPoint
	Location postgis.PointS `db:"location" json:"-"`
}

type AddGeoPoint struct {
	Title           string    `json:"title" example:"Forêt à l'aube" validate:"required,min=3,max=30"`
	UserId          int       `json:"userId" example:"1" validate:"isdefault"`
	Latitude        float64   `json:"latitude" example:"38.652608" validate:"required,latitude"`
	Longitude       float64   `json:"longitude" example:"-120.357448" validate:"required,longitude"`
	Date            time.Time `json:"date" example:"2022-05-26T11:17:35.079344Z" validate:"required,lt"`
	Amplitudes      []float64 `json:"amplitudes" example:"0,1,2,3,45,3,2,1" validate:"required,min=100,max=1000"`
	PictureTemplate string    `json:"picture_template" example:"forest" validate:"omitempty,oneof=forest sea mountain swamp"`
}

type BindGeoPoint struct {
	Sound   *multipart.FileHeader `form:"sound" binding:"required"`
	Picture *multipart.FileHeader `form:"picture" binding:"omitempty"`
}

type ClosestGeoPoint struct {
	Latitude   float64       `uri:"latitude" example:"18.16255" validate:"latitude"`
	Longitude  float64       `uri:"longitude" example:"40.35735" validate:"longitude"`
	SRID       *int32        `form:"srid" example:"4326" validate:"omitempty"`
	IdExcluded pq.Int32Array `form:"not[]" example:"1,2,3,4" validate:"lt=10"`
}

type ClosestGeoId struct {
	Id int `json:"id" db:"id" example:"18"`
}

const WGS84 = 4326
