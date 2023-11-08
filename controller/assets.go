package controller

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haran/biophonie-api/database"
)

const MINSIZE = 49

func (c *Controller) refreshGeoJson() {
	geoJson := make([]byte, 0)
	err := c.Db.Get(&geoJson, database.GeosAsGeoJson)
	if err != nil {
		log.Fatalf("refreshGeoJson: could not query geojson: %s", err)
	}

	file, err := os.OpenFile(c.geoJsonPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("refreshGeoJson: could not open geojson file: %s", err)
	}
	defer file.Close()

	if _, err := file.Write(geoJson); err != nil {
		log.Fatalf("refreshGeoJson: could not write geojsonfile: %s", err)
	}
}

func (c *Controller) ClearGeoPoint(ctx *gin.Context) {
	ctx.Next()
	picture := ctx.GetString("picture")
	sound := ctx.GetString("sound")

	if err := os.Remove(picture); err != nil {
		log.Println("could not rm picture: ", err)
	}
	if err := os.Remove(sound); err != nil {
		log.Println("could not rm sound: ", err)
	}

	c.refreshGeoJson()
}

func (c *Controller) AppendGeoJson(ctx *gin.Context) {
	ctx.Next()
	geoId := ctx.GetUint64("geoId")

	featBytes := make([]byte, 0)
	err := c.Db.Get(&featBytes, database.GeoAsFeat, geoId)
	if err != nil {
		log.Println("could not get geopoint to append geojson: ", err)
		return
	}

	f, err := os.OpenFile(c.geoJsonPath, os.O_RDWR, 0600)
	if err != nil {
		log.Println("could not open geojson to append feature: ", err)
		return
	}
	defer f.Close()

	info, err := os.Stat(f.Name())
	if err != nil {
		log.Println("could not stat GeoJson file: ", err)
		return
	}

	offset := info.Size()
	fmt.Println(offset)
	if offset == MINSIZE {
		featBytes = append([]byte("["), featBytes...)
		featBytes = append(featBytes, []byte("]}")...)
		offset -= 5
	} else {
		featBytes = append([]byte(","), featBytes...)
		featBytes = append(featBytes, []byte("]}")...)
		offset -= 2
	}
	if n, err := f.WriteAt(featBytes, offset); err != nil || n != len(featBytes) {
		log.Println("could not refresh GeoJson file: ", err)
		return
	}
}
