package controller

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter(c *Controller) *gin.Engine {
	r := gin.Default()
	r.Use(c.HandleErrors)
	r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = 10000000 // 10 MB
	v1 := r.Group("/api/v1")
	{
		v1.Static("/assets", "./public")
		users := v1.Group("/user")
		{
			users.GET("/:name", c.GetUser)
			users.POST("", c.CreateUser)
		}
		geopoints := v1.Group("/geopoint")
		{
			geopoints.GET("/:id", c.GetGeoPoint)
			geopoints.POST("", c.BindGeoPoint, c.CheckGeoFiles, c.CreateGeoPoint)
			geopoints.GET("/:id/picture", c.GetPicture)
			geopoints.GET("/:id/sound", c.GetSound)
		}
	}

	// r.UseH2C = true // try to use http2 maybe with next version of gin-gonic
	r.GET("/ping", c.Pong)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}