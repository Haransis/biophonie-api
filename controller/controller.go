package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Param user body user.AddUser true "Add user"
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

	if _, err := c.Db.Exec("INSERT INTO accounts (username, created_on) VALUES ($1,now())", addUser.Name); err != nil {
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

	ctx.JSON(http.StatusCreated, addUser)
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
// @Router /user/{username} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")

	var user user.User
	if err := c.Db.Get(&user, "SELECT * FROM accounts WHERE name = $1", name); err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound,
				fmt.Sprintf("error reading accounts: unknown user %s", name))
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Panicf("error reading accounts: %q", err)
		}
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /geopoint/{id} [get]
func (c *Controller) GetGeoPoint(ctx *gin.Context) {
	id := ctx.Param("id")

	var geopoint geopoint.GeoPoint
	if err := c.Db.Get(&geopoint, "SELECT * FROM geopoints WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound,
				fmt.Sprintf("error reading geopoints: unknown geopoint %s", id))
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Panicf("error reading geopoints: %q", err)
		}
	}

	ctx.JSON(http.StatusOK, geopoint)
}

// CreateGeoPoint godoc
// @Summary create a geopoint
// @Description create the geopoint in the database and receive the sound and picture file
// @Accept mpfd
// @Produce json
// @Tags Geopoint
// @Param geopoint formData geopoint.AddGeoPoint true "geopoint infos"
// @Param sound formData file true "geopoint sound"
// @Param picture formData file true "geopoint picture"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /geopoint [post]
func (c *Controller) CreateGeoPoint(ctx *gin.Context) {
	// TODO (not implemented)

	ctx.JSON(http.StatusOK, "OK")
}

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
