package api

import (
	"github.com/gin-gonic/gin"
)

// Listen starts the api server at the specified address
func Listen(addr string) error {
	r := gin.Default()

	api := r.Group("/api/v1")

	api.POST("scan/:studentID", scanHandler)

	return r.Run(addr)
}