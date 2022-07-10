package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haran/biophonie-api/controller/geopoint"
)

const MINSIZE = 43

func (c *Controller) refreshGeoJson() {
	geos := make([]geopoint.GeoPoint, 0)
	err := c.Db.Select(&geos, "SELECT * FROM geopoints WHERE available=TRUE")
	if err != nil {
		log.Fatalf("refreshGeoJson: could not query geopoints: %s", err)
	}

	geoJson := geopoint.ToGeoJson(geos)

	file, err := os.OpenFile(c.geoJsonPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("refreshGeoJson: could not open geojson file: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(geoJson); err != nil {
		log.Fatalf("refreshGeoJson: could not write geojson: %s", err)
	}
}

func (c *Controller) AppendGeoJson(ctx *gin.Context) {
	ctx.Next()
	geoId := ctx.GetUint64("geoId")

	var geopoint geopoint.GeoPoint
	err := c.Db.Get(&geopoint, "SELECT * FROM geopoints WHERE id=$1", geoId)
	if err != nil {
		log.Println("could not get geopoint to append geojson: ", err)
		return
	}

	feat := geopoint.ToFeature()
	featBytes, err := json.Marshal(feat)
	if err != nil {
		log.Fatalf("could not marshal feature: %s", err)
	}

	f, err := os.OpenFile(c.geoJsonPath, os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	info, err := os.Stat(f.Name())
	if err != nil {
		log.Fatalf("could not stat GeoJson file: %s", err)
	}

	offset := info.Size()
	fmt.Println(offset)
	if offset == MINSIZE {
		featBytes = append(featBytes, []byte("]}")...)
		offset -= 3
	} else {
		featBytes = append([]byte(","), featBytes...)
		featBytes = append(featBytes, []byte("]}")...)
		offset -= 2
	}
	if n, err := f.WriteAt(featBytes, offset); err != nil || n != len(featBytes) {
		log.Fatalf("could not refresh GeoJson file: %s", err)
	}
}
