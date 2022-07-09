package geopoint

type GeoJson struct {
	Type     string    `json:"type" default:"FeatureCollection"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type        string     `json:"type" default:"Point"`
	Coordinates []float64  `json:"coordinates"`
	Properties  Properties `json:"properties"`
}

type Properties struct {
	Name  string `json:"name"`
	Id    int    `json:"id"`
	Cache bool   `json:"cache" default:"false"`
}
