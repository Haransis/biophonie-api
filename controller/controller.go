package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	if _, err := c.Db.Exec("INSERT INTO accounts (username, created_on, last_login) VALUES ($1,now(),now())", addUser.Username); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				httputil.NewError(ctx, http.StatusConflict, fmt.Errorf("user with username %s already exists", addUser.Username))
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
// @Param username path string true "user name"
// @Success 200 {object} user.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/{username} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	username := ctx.Param("username")

	var user user.User
	if err := c.Db.Get(&user, "SELECT * FROM accounts WHERE username = $1", username); err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusInternalServerError,
				fmt.Sprintf("error reading accounts: unknown username %s", username))
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Panicf("error while reading accounts: %q", err)
		}
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) Pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
