package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cridenour/go-postgis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haran/biophonie-api/controller/geopoint"
	"github.com/haran/biophonie-api/controller/user"
	"github.com/haran/biophonie-api/database"
	"github.com/haran/biophonie-api/httputil"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Controller struct {
	Db *sqlx.DB
}

func NewController() *Controller {
	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("Error initializing database: %q", err)
	}

	return &Controller{Db: db}
}

// CreateUser godoc
// @Summary create user
// @Description create a user in the database
// @Accept json
// @Produce json
// @Tags User
// @Param user body user.User true "Add user"
// @Success 201 {object} user.AddUser
// @Failure 400 {object} httputil.HTTPError
// @Failure 409 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var addUser user.AddUser
	if err := ctx.ShouldBindJSON(&addUser); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	stmt, err := c.Db.PrepareNamed("INSERT INTO accounts (name, created_on) VALUES (:name,now()) RETURNING id")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Panicf("error creating database request: %q", err)
	}

	var id int
	if err := stmt.Get(&id, addUser); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				httputil.NewError(ctx, http.StatusConflict, fmt.Errorf("user with username %s already exists", addUser.Name))
				return
			} else {
				httputil.NewError(ctx, http.StatusInternalServerError, errors.New(err.Code.Name()))
				return
			}
		}
		log.Panic(err)
	}

	var user user.User
	c.GetElementById(ctx, id, "users", user)
}

// GetUser godoc
// @Summary get a user
// @Description retrieve the user in the database using its name
// @Accept json
// @Produce json
// @Tags User
// @Param name path string true "user name"
// @Success 200 {object} user.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/{name} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")

	var user user.User
	c.GetUserByName(ctx, name, &user)
}

// GetGeoPoint godoc
// @Summary get a geopoint
// @Description retrieve the geopoint in the database using its name
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /geopoint/{id} [get]
func (c *Controller) GetGeoPoint(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
	}

	var geopoint geopoint.GeoPoint
	c.GetElementById(ctx, id, "geopoints", &geopoint)
}

// CreateGeoPoint godoc
// @Summary create a geopoint
// @Description create the geopoint in the database and receive the sound and picture file (see testgeopoint dir)
// @Accept mpfd
// @Produce json
// @Tags Geopoint
// @Param geopoint formData file true "geopoint infos in a utf-8 json file"
// @Param sound formData file true "geopoint sound"
// @Param picture formData file true "geopoint picture"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /geopoint [post]
func (c *Controller) CreateGeoPoint(ctx *gin.Context) {
	var bindGeo geopoint.BindGeopoint
	if err := ctx.ShouldBind(&bindGeo); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if !httputil.CheckFileContentType(bindGeo.Sound, "audio/wave") {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("sound was not wave file"))
		return
	}

	if !httputil.CheckFileContentType(bindGeo.Picture, "image/jpeg") {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("image was not jpeg file"))
		return
	}

	geoFile, err := bindGeo.Geopoint.Open()
	if err != nil {
		log.Panicln(err)
	}
	defer geoFile.Close()

	geoBytes, err := ioutil.ReadAll(geoFile)
	if err != nil {
		log.Panicln(err)
	}

	var addGeoPoint geopoint.AddGeoPoint
	if err := json.Unmarshal(geoBytes, &addGeoPoint); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("could not parse geopoint file"))
		return
	}

	var userExists bool
	if err := c.Db.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id=$1)", addGeoPoint.UserId); err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Panicf("error reading geopoint: %q", err)
	}

	if !userExists {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("user was not found"))
		return
	}

	geoPoint := geopoint.GeoPoint{
		Title:  addGeoPoint.Title,
		UserId: addGeoPoint.UserId,
		Location: postgis.Point{
			X: addGeoPoint.Location[0],
			Y: addGeoPoint.Location[1],
		},
		CreatedOn:  addGeoPoint.Date,
		Amplitudes: addGeoPoint.Amplitudes,
		Picture:    fmt.Sprintf("%s.jpg", uuid.New()),
		Sound:      fmt.Sprintf("%s.wav", uuid.New()),
	}

	stmt, err := c.Db.PrepareNamed("INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on) VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on) RETURNING id")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Panicf("error creating database request: %q", err)
	}

	if err := stmt.Get(&geoPoint.Id, geoPoint); err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Panicf("error creating geopoint: %q", err)
	}

	if err := ctx.SaveUploadedFile(bindGeo.Picture, fmt.Sprintf("./public/picture/%s", geoPoint.Picture)); err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("could not save picture file"))
		log.Panicf("could not save picture file: %s", err)
	}

	if err := ctx.SaveUploadedFile(bindGeo.Sound, fmt.Sprintf("./public/sound/%s", geoPoint.Sound)); err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("could not save sound file"))
		log.Panicf("could not save sound file: %s", err)
	}

	ctx.JSON(http.StatusOK, geoPoint)
}

// TODO replace httputil.NewError by ctx.Abort() ?

// GetPicture godoc
// @Summary get the url of the picture
// @Description create the geopoint in the database and receive the sound and picture file
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /geopoint/{id}/picture [get]
func (c *Controller) GetPicture(ctx *gin.Context) {
	// TODO (not implemented)

	ctx.JSON(http.StatusOK, "OK")
}

// GetPicture godoc
// @Summary pings the api
// @Description used to check if api is alive
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 500 {object} httputil.HTTPError
// @Router /ping [get]
func (c *Controller) Pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
