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

	api.POST("scan", controllers.scanHandler)

	return r.Run(addr)
}
