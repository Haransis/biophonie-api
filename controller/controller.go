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
)

type Controller struct {
	Db *sqlx.DB
}

func NewController() *Controller {
	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("error initializing database: %q", err)
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
// @Failure 400 {object} controller.ErrMsg
// @Failure 409 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /user [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var addUser user.AddUser
	if err := ctx.BindJSON(&addUser); err != nil {
		return
	}

	stmt, err := c.Db.PrepareNamed("INSERT INTO accounts (name, created_on) VALUES (:name,now()) RETURNING id")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not prepare user creation: %s", err))
		return
	}

	var id int
	if err := stmt.Get(&id, addUser); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not create user")
		ctx.Abort()
		return
	}

	var user user.User
	if err := c.Db.Get(&user, "SELECT * FROM accounts WHERE id = $1", id); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not retrieve created user")
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetUser godoc
// @Summary get a user
// @Description retrieve the user in the database using its name
// @Accept json
// @Produce json
// @Tags User
// @Param name path string true "user name"
// @Success 200 {object} user.User
// @Failure 400 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /user/{name} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")

	var user user.User
	if err := c.Db.Get(&user, "SELECT * FROM accounts WHERE name = $1", name); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get user")
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetGeoPoint godoc
// @Summary get a geopoint
// @Description retrieve the geopoint in the database using its name
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 400 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /geopoint/{id} [get]
func (c *Controller) GetGeoPoint(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	var geopoint geopoint.GeoPoint
	if err := c.Db.Get(&geopoint, "SELECT * FROM geopoints WHERE id = $1", id); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get geopoint")
		ctx.Abort()
		return
	}
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
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /geopoint [post]
func (c *Controller) CreateGeoPoint(ctx *gin.Context) {
	var bindGeo geopoint.BindGeopoint
	if err := ctx.Bind(&bindGeo); err != nil {
		return
	}

	if !httputil.CheckFileContentType(bindGeo.Sound, "audio/wave") {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("sound was not wave file")).SetType(gin.ErrorTypePublic)
		return
	}

	if !httputil.CheckFileContentType(bindGeo.Picture, "image/jpeg") {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("image was not jpeg file")).SetType(gin.ErrorTypePublic)
		return
	}

	geoFile, err := bindGeo.Geopoint.Open()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not open geofile: %s", err))
	}
	defer geoFile.Close()

	geoBytes, err := ioutil.ReadAll(geoFile)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not read geofile: %s", err))
	}

	var addGeoPoint geopoint.AddGeoPoint
	if err := json.Unmarshal(geoBytes, &addGeoPoint); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	var userExists bool
	if err := c.Db.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id=$1)", addGeoPoint.UserId); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not check if user exists: %s", err))
		return
	}

	if !userExists {
		ctx.AbortWithError(http.StatusNotFound, errors.New("user was not found")).SetType(gin.ErrorTypePublic)
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
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not prepare geopoint creation: %s", err))
		return
	}

	if err := stmt.Get(&geoPoint.Id, geoPoint); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not create geopoint")
		ctx.Abort()
		return
	}

	if err := ctx.SaveUploadedFile(bindGeo.Picture, fmt.Sprintf("./public/picture/%s", geoPoint.Picture)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not save uploaded picture: %s", err))
		return
	}

	if err := ctx.SaveUploadedFile(bindGeo.Sound, fmt.Sprintf("./public/sound/%s", geoPoint.Sound)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not save uploaded sound: %s", err))
		return
	}

	ctx.JSON(http.StatusOK, geoPoint)
}

// GetPicture godoc
// @Summary get the url of the picture
// @Description create the geopoint in the database and receive the sound and picture file
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {string} string
// @Failure 400 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
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
// @Failure 500 {object} controller.ErrMsg
// @Router /ping [get]
func (c *Controller) Pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
