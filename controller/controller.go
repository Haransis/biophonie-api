package controller

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/haran/biophonie-api/database"
	"github.com/jmoiron/sqlx"
)

const geoJsonFileName = "geojson.json"

type Controller struct {
	Db           *sqlx.DB
	assetsFolder string
	webFolder    string
	geoJsonPath  string
	verifyKey    *rsa.PublicKey
	signKey      *rsa.PrivateKey
}

func NewController() *Controller {
	c := &Controller{}
	c.readKeys()

	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("error initializing database: %q", err)
	}
	c.Db = db

	c.assetsFolder = os.Getenv("ASSETS_FOLDER")
	if c.assetsFolder == "" {
		log.Fatalf("assets path is empty")
	}

	c.webFolder = os.Getenv("WEB_FOLDER")
	if c.webFolder == "" {
		log.Fatalf("web path is empty")
	}

	c.geoJsonPath = c.assetsFolder + string(os.PathSeparator) + geoJsonFileName
	c.refreshGeoJson()

	return c
}
