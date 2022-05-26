package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/haran/biophonie-api/controller"
	_ "github.com/haran/biophonie-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger for biophonie-api
// @version         1.0
// @description     API of biophonie (https://secret-garden-77375.herokuapp.com/).
// @termsOfService  TODO

// @contact.name   TODO
// @contact.url    TODO
// @contact.email  TODO

// @license.name  GPL-3.0 license
// @license.url   https://www.gnu.org/licenses/gpl-3.0.en.html

// @BasePath /api/v1

func main() {

	controller := controller.NewController()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/user")
		{
			users.GET("/:username", controller.GetUser)
			users.POST("", controller.CreateUser)
		}
	}
	r.GET("/ping", controller.Pong)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Stopping server: %q", err)
	}
}
