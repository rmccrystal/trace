package api

import "github.com/gin-gonic/gin"

// /api/v1/scan/:studentID
// Called whenever someone scans their barcode
func scanHandler(c *gin.Context) {
	studentID := c.Param("studentID")
}
