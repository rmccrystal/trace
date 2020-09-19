package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"trace/pkg/controllers"
	"trace/pkg/database"
)

// Listen starts the api server at the specified address
func Listen(addr string, config *Config) error {
	if database.DB == nil {
		return errors.New("database is not connected")
	}

	GlobalConfig = config

	r := gin.Default()

	api := r.Group("/api/v1")

	api.POST("scan", controllers.OnScan)

	api.POST("location/create", controllers.CreateLocation)
	api.GET("location", controllers.GetLocations)
	api.GET("location/:id", controllers.GetLocationByID)
	api.DELETE("location/:id", controllers.DeleteLocation)
	api.PATCH("location/:id", controllers.UpdateLocation)

	return r.Run(addr)
}
