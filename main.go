package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/haran/biophonie-api/database"
	"github.com/haran/biophonie-api/router"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("Error initializing database: %q", err)
	}

	router := &router.Router{
		Db: db,
	}

	r := gin.Default()
	r.GET("/ping", router.Pong())
	r.GET("/user", router.GetUser())
	r.POST("/user", router.CreateUser())

	err = r.Run()
	if err != nil {
		log.Fatalf("Stopping server: %q", err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
