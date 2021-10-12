package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		fmt.Println("stopping server: ", err)
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}