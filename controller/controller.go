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
	Db          *sqlx.DB
	publicPath  string
	geoJsonPath string
	verifyKey   *rsa.PublicKey
	signKey     *rsa.PrivateKey
}

func NewController() *Controller {
	c := &Controller{}
	c.readKeys()

	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("error initializing database: %q", err)
	}
	c.Db = db

	c.publicPath = os.Getenv("PUBLIC_PATH")
	if c.publicPath == "" {
		log.Fatalf("public path is empty")
	}
	c.geoJsonPath = c.publicPath + string(os.PathSeparator) + geoJsonFileName
	c.refreshGeoJson()

	return c
}
