package controller

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter(c *Controller) *gin.Engine {
	r := gin.Default()
	r.Use(c.HandleErrors)
	r.Use(gin.Recovery())
	// r.Use(cors.Default()) // use for sound streaming
	r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = 100000000 // 100 MB
	r.Use(static.Serve("/", static.LocalFile(c.webFolder, false)))
	v1 := r.Group("/api/v1")
	{
		v1.Static("/assets", c.assetsFolder)
		users := v1.Group("/user")
		{
			users.GET("/:name", c.GetUser)
			users.POST("", c.PostUser)
			users.POST("/authorize", c.AuthorizeUser)
		}
		geopoints := v1.Group("/geopoint")
		{
			geopoints.GET("/:id", c.GetGeoPoint)
			geopoints.GET("/closest/to/:latitude/:longitude", c.GetClosestGeoPoint)
			geopoints.GET("/:id/assets", c.GetAssets)
		}
		restricted := v1.Group("/restricted", c.Authorize)
		{
			restricted.POST("/geopoint", c.BindGeoPoint, c.CreateGeoPoint)
			restricted.GET("/ping", c.AuthPong)
			toAdmins := restricted.Group("", c.AuthorizeAdmin)
			{
				toAdmins.PATCH("/geopoint/:id/enable", c.EnableGeoPoint, c.AppendGeoJson)
				toAdmins.PATCH("/user/:id", c.MakeAdmin)
				toAdmins.GET("/geopoint/:id", c.GetGeoPoint)
				toAdmins.DELETE("/geopoint/:id", c.DeleteGeoPoint, c.ClearGeoPoint)
			}
		}
		v1.GET("/ping", c.Pong)
	}

	// r.UseH2C = true // try to use http2 maybe with next version of gin-gonic
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
