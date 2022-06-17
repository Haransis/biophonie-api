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
// @description     API of biophonie (https://secret-garden-77375.herokuapp.com/). Files are located in "assets/"
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
	r.MaxMultipartMemory = 10000000 // 10 MB
	v1 := r.Group("/api/v1")
	{
		v1.Static("/assets", "./public")
		users := v1.Group("/user")
		{
			users.GET("/:name", controller.GetUser)
			users.POST("", controller.CreateUser)
		}
		geopoints := v1.Group("/geopoint")
		{
			geopoints.GET("/:id", controller.GetGeoPoint)
			geopoints.POST("", controller.CreateGeoPoint)
			geopoints.GET("/:id/picture", controller.GetPicture)
		}
	}

	// r.UseH2C = true // try to use http2 maybe with next version of gin-gonic
	r.GET("/ping", controller.Pong)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Stopping server: %q", err)
	}
}
