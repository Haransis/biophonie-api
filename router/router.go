package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Db *sql.DB
}

func (r *Router) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.String(http.StatusBadRequest, "username cannot be null")
			return
		}

		if _, err := r.Db.Exec("INSERT INTO accounts VALUES ($1,,,now())", username); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating user: %q", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "created",
			"username": username,
		})
	}
}

func (r *Router) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("username")

		var email string
		if err := r.Db.QueryRow("SELECT email FROM accounts WHERE username = '$1'", name).Scan(&email); err != nil {
			if err == sql.ErrNoRows {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading accounts: unknown username %s", name))
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"email": name,
		})
	}
}

func (r *Router) Pong() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
