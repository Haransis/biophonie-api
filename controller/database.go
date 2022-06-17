package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haran/biophonie-api/controller/user"
)

func (c *Controller) GetUserByName(ctx *gin.Context, name string, user *user.User) {
	if err := c.Db.Get(user, "SELECT * FROM accounts WHERE name = $1", name); err != nil {
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

func (c *Controller) GetElementById(ctx *gin.Context, id int, table string, element interface{}) {
	if err := c.Db.Get(element, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table), id); err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound,
				fmt.Sprintf("error reading %s: unknown element %d", table, id))
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			log.Panicf("error reading %s: %q", table, err)
		}
	}
	ctx.JSON(http.StatusOK, element)
}
