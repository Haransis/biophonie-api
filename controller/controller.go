package controller

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cridenour/go-postgis"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/haran/biophonie-api/controller/geopoint"
	"github.com/haran/biophonie-api/controller/user"
	"github.com/haran/biophonie-api/database"
	"github.com/haran/biophonie-api/httputil"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const geoJsonFileName = "/geojson.json"

const requestClosestGeoId = `--sql
	WITH excluded(id) AS ( SELECT UNNEST($2::int[])) 
	SELECT geo.id FROM geopoints geo 
	WHERE NOT EXISTS(SELECT 1 FROM excluded e WHERE geo.id = e.id) AND available = TRUE
	ORDER BY geo.location <-> GeomFromEWKB($1)
	LIMIT 1;
`

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
	c.geoJsonPath = c.publicPath + geoJsonFileName
	c.refreshGeoJson()

	return c
}

// CreateUser godoc
// @Summary create user
// @Description create a user in the database
// @Accept json
// @Produce json
// @Tags User
// @Param user body user.AddUser true "Add user"
// @Success 200 {object} user.User
// @Failure 400 {object} controller.ErrMsg
// @Failure 409 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /user [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var addUser user.AddUser
	if err := ctx.BindJSON(&addUser); err != nil {
		return
	}

	// generate password and hash it
	password := uuid.New().String()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not hash password: %s", err))
		return
	}

	stmt, err := c.Db.Preparex("INSERT INTO accounts (name, created_on, password) VALUES ($1,now(),$2) RETURNING id")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not prepare user creation: %s", err))
		return
	}

	var id int
	if err := stmt.Get(&id, addUser.Name, hashedToken); err != nil {
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
	user.Password = password

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
	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

// GetRestrictedGeoPoint godoc
// @Summary get a geopoint which was not enabled
// @Description retrieve the geopoint in the database using its id
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 400 {object} controller.ErrMsg
// @Failure 403 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/geopoint/{id} [get]
func GetRestrictedGeoPoint() {}

// GetGeoPoint godoc
// @Summary get a geopoint
// @Description retrieve the geopoint in the database using its id
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 400 {object} controller.ErrMsg
// @Failure 403 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /geopoint/{id} [get]
func (c *Controller) GetGeoPoint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	var geopoint geopoint.DbGeoPoint
	if err := c.Db.Get(&geopoint, "SELECT * FROM geopoints WHERE id = $1", id); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get geopoint")
		ctx.Abort()
		return
	}

	if !geopoint.Available && !ctx.GetBool("admin") {
		ctx.AbortWithError(http.StatusForbidden, errors.New("geopoint is not enabled yet")).SetType(gin.ErrorTypePublic)
		return
	}
	geopoint.Latitude = geopoint.Location.Y
	geopoint.Longitude = geopoint.Location.X

	ctx.JSON(http.StatusOK, geopoint)
}

// GetClosestGeoPoint godoc
// @Summary get the closest geopoint
// @Description retrieve the closest geopoint to another geopoint, excluding other geopoints
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param latitude path float64 true "latitude"
// @Param longitude path float64 true "longitude"
// @Param srid query int32 false "srid to project"
// @Param not[] query []int32 false "optional ids to exclude from search"
// @Success 200 {object} geopoint.ClosestGeoId
// @Failure 400 {object} controller.ErrMsg
// @Failure 403 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /geopoint/closest/to/{latitude}/{longitude} [get]
func (c *Controller) GetClosestGeoPoint(ctx *gin.Context) {
	var closestTo geopoint.ClosestGeoPoint
	if err := ctx.BindUri(&closestTo); err != nil {
		return
	}
	if err := ctx.BindQuery(&closestTo); err != nil {
		return
	}

	target := postgis.PointS{X: closestTo.Latitude, Y: closestTo.Longitude, SRID: geopoint.WGS84}
	if closestTo.SRID != nil {
		target.SRID = *closestTo.SRID
	}

	var geoId geopoint.ClosestGeoId
	if err := c.Db.Get(&geoId, requestClosestGeoId, target, closestTo.IdExcluded); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get closest geopoint")
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, geoId)
}

// DeleteGeopoint godoc
// @Summary delete a geopoint
// @Description delete a geopoint
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} string "user is now admin"
// @Failure 400 {object} controller.ErrMsg
// @Failure 401 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/geopoint/{id} [delete]
func (c *Controller) DeleteGeoPoint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := c.Db.Exec("DELETE FROM geopoints WHERE id = $1", id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected != 1 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("not found or already deleted"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "geopoint was deleted"})
}

// AuthorizeUser godoc
// @Summary create a token
// @Description create a token
// @Accept json
// @Produce json
// @Tags Authentication
// @Param user body user.AuthUser true "authentication user"
// @Success 200 {string} user.AccessToken "token to use for authentication"
// @Failure 400 {object} controller.ErrMsg
// @Failure 401 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /user/authorize [post]
func (c *Controller) AuthorizeUser(ctx *gin.Context) {
	var authUser user.AuthUser
	if err := ctx.BindJSON(&authUser); err != nil {
		return
	}

	var authorizedUser user.User
	if err := c.Db.Get(&authorizedUser, "SELECT * FROM accounts WHERE name = $1", authUser.Name); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get password")
		ctx.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authorizedUser.Password), []byte(authUser.Password)); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not compare password and hash")
		ctx.Abort()
		return
	}

	token, err := c.createToken(authorizedUser.Name, authorizedUser.Admin)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not sign token: %s", err))
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, user.AccessToken{Token: token})
}

// MakeAdmin godoc
// @Summary make a user admin
// @Description make a user admin
// @Accept json
// @Produce plain
// @Tags Authentication
// @Param id path int true "user id"
// @Success 200 {string} string "user is now admin"
// @Param Authorization header string true "Authentication header"
// @Failure 400 {object} controller.ErrMsg
// @Failure 403 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/user/{id} [patch]
func (c *Controller) MakeAdmin(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := c.Db.Exec("UPDATE accounts SET admin = TRUE WHERE id = $1", id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected != 1 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("not found")).SetType(gin.ErrorTypePublic)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user is now admin"})
}

// BindGeoPoint godoc
// @Summary create a geopoint
// @Description create the geopoint in the database and save the sound and picture file (see testgeopoint dir)
// @Accept mpfd
// @Produce json
// @Tags Geopoint
// @Param geopoint formData string true "geopoint infos in a utf-8 json file"
// @Param sound formData file true "geopoint sound"
// @Param picture formData file false "geopoint picture"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} geopoint.GeoPoint
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/geopoint [post]
func (c *Controller) BindGeoPoint(ctx *gin.Context) {
	var bindGeo geopoint.BindGeoPoint
	if err := ctx.Bind(&bindGeo); err != nil {
		return
	}

	var geoPointJson string
	if geoPointJson = ctx.Request.PostFormValue("geopoint"); geoPointJson == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("geopoint was not in request")).
			SetType(gin.ErrorTypeBind)
		return
	}

	var addGeo geopoint.AddGeoPoint
	if err := json.Unmarshal([]byte(geoPointJson), &addGeo); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	if err := validator.New().Struct(addGeo); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypeBind)
		return
	}

	soundName := uuid.NewString() + ".wav"
	pictureName := uuid.NewString() + ".jpg"
	if !httputil.CheckFileContentType(bindGeo.Sound, "audio/wave") {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("sound was not wave file")).SetType(gin.ErrorTypePublic)
		return
	}

	if addGeo.PictureTemplate == "" {
		if !httputil.CheckFileContentType(bindGeo.Picture, "image/jpeg") {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("image was not jpeg file")).SetType(gin.ErrorTypePublic)
			return
		}
	} else {
		pictureName = addGeo.PictureTemplate + ".jpg"
	}

	addGeo.UserId, _ = ctx.MustGet("userId").(int)

	geoPoint := geopoint.GeoPoint{
		Title:      addGeo.Title,
		UserId:     addGeo.UserId,
		Latitude:   addGeo.Latitude,
		Longitude:  addGeo.Longitude,
		CreatedOn:  addGeo.Date,
		Amplitudes: addGeo.Amplitudes,
		Picture:    pictureName,
		Sound:      soundName,
	}

	ctx.Set("bindGeo", bindGeo)
	ctx.Set("geoPoint", geoPoint)
}

func (c *Controller) CreateGeoPoint(ctx *gin.Context) {
	bindGeo, _ := ctx.MustGet("bindGeo").(geopoint.BindGeoPoint)
	geoPoint, _ := ctx.MustGet("geoPoint").(geopoint.GeoPoint)

	dbGeoPoint := geopoint.DbGeoPoint{GeoPoint: &geoPoint, Location: postgis.PointS{SRID: geopoint.WGS84, X: geoPoint.Longitude, Y: geoPoint.Latitude}}
	stmt, err := c.Db.PrepareNamed("INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on) VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on) RETURNING id")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not prepare geopoint creation: %s", err))
		return
	}

	if err := stmt.Get(&geoPoint.Id, dbGeoPoint); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not create geopoint")
		ctx.Abort()
		return
	}

	if bindGeo.Picture != nil {
		if err := ctx.SaveUploadedFile(bindGeo.Picture, fmt.Sprintf("%s/picture/%s", c.publicPath, geoPoint.Picture)); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not save uploaded picture: %s", err))
			return
		}
	}

	if err := ctx.SaveUploadedFile(bindGeo.Sound, fmt.Sprintf("%s/sound/%s", c.publicPath, geoPoint.Sound)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not save uploaded sound: %s", err))
		return
	}

	ctx.JSON(http.StatusOK, geoPoint)
}

// EnableGeoPoint godoc
// @Summary make the geopoint available
// @Description make the geopoint available
// @Accept json
// @Produce plain
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} string
// @Failure 400 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/geopoint/{id}/enable [patch]
func (c *Controller) EnableGeoPoint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := c.Db.Exec("UPDATE geopoints SET available = TRUE WHERE id = $1 AND available = FALSE", id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected != 1 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("not found or already enabled"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "geopoint was enabled"})
	ctx.Set("geoId", id)
}

// GetPicture godoc
// @Summary get the picture filename
// @Description located in assets/
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	var picture string
	if err := c.Db.Get(&picture, "SELECT picture FROM geopoints WHERE id = $1", id); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get picture")
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, picture)
}

// GetSound godoc
// @Summary get the sound filename
// @Description located in assets/
// @Accept json
// @Produce json
// @Tags Geopoint
// @Param id path int true "geopoint id"
// @Success 200 {string} string
// @Failure 400 {object} controller.ErrMsg
// @Failure 404 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /geopoint/{id}/sound [get]
func (c *Controller) GetSound(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	var sound string
	if err := c.Db.Get(&sound, "SELECT sound FROM geopoints WHERE id = $1", id); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get sound")
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, sound)
}

// AuthPong godoc
// @Summary pings the authenticated api
// @Description used to check if client is authenticated
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} string
// @Failure 403 {object} controller.ErrMsg
// @Failure 500 {object} controller.ErrMsg
// @Router /restricted/ping [get]
func (c *Controller) AuthPong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "authenticated pong",
	})
}

// Pong godoc
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
