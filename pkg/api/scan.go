package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// /api/v1/scan/:studentID
// Called whenever someone scans their barcode
func scanHandler(c *gin.Context) {
	scanRequest := struct {
		StudentID string
		LocationID string
	}{}

	if err := c.BindJSON(&scanRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, )
	}
}
