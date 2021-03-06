package geopoint

import (
	"mime/multipart"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/lib/pq"
)

type GeoPoint struct {
	Id         int           `db:"id" json:"id" example:"1" validate:"gt=0"`
	Title      string        `db:"title" json:"title" example:"Forêt à l'aube"`
	UserId     int           `db:"user_id" json:"userId" example:"1" validate:"gt=0"`
	Location   postgis.Point `db:"location" json:"location"`
	CreatedOn  time.Time     `db:"created_on" json:"createdOn" example:"2022-05-26T11:17:35.079344Z"`
	Amplitudes pq.Int64Array `db:"amplitudes" json:"amplitudes" swaggertype:"array,number" example:"0,1,2,3,45,3,2,1"`
	Picture    string        `db:"picture" json:"picture" example:"https://example.com/picture-1.jpg"`
	Sound      string        `db:"sound" json:"sound" example:"https://example.com/sound-2.wav"`
	Available  bool          `db:"available" json:"available" example:"true"`
}

type AddGeoPoint struct {
	Title           string    `json:"title" example:"Forêt à l'aube" validate:"required,min=3,max=30"`
	UserId          int       `json:"userId" example:"1" validate:"isdefault"`
	Latitude        float64   `json:"latitude" example:"38.652608" validate:"required,latitude"`
	Longitude       float64   `json:"longitude" example:"-120.357448" validate:"required,longitude"`
	Date            time.Time `json:"date" example:"2022-05-26T11:17:35.079344Z" validate:"required,lt=(time.Time)"`
	Amplitudes      []int64   `json:"amplitudes" example:"0,1,2,3,45,3,2,1" validate:"required,min=10,max=500"`
	PictureTemplate string    `json:"picture_template" example:"forest" validate:"omitempty,oneof=forest sea mountain swamp"`
}

type BindGeoPoint struct {
	Geopoint *multipart.FileHeader `form:"geopoint" binding:"required"`
	Sound    *multipart.FileHeader `form:"sound" binding:"required"`
	Picture  *multipart.FileHeader `form:"picture"`
}

func ToGeoJson(gs []GeoPoint) *GeoJson {
	if len(gs) == 0 {
		return &GeoJson{
			Type:     "FeatureCollection",
			Features: make([]Feature, 0),
		}
	}
	feats := make([]Feature, 0)
	for _, geopoint := range gs {
		feats = append(feats, geopoint.ToFeature())
	}
	return &GeoJson{
		Type:     "FeatureCollection",
		Features: feats,
	}
}

func (g *GeoPoint) ToFeature() Feature {
	return Feature{
		Type:        "Point",
		Coordinates: []float64{g.Location.X, g.Location.Y},
		Properties: Properties{
			Name: g.Title,
			Id:   g.Id,
		},
	}
}
