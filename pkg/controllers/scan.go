package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"trace/pkg/trace"
)

// POST /api/v1/scan
// Called whenever someone scans their barcode
func OnScan(c *gin.Context) {
	scanRequest := struct {
		StudentHandle string `json:"student_handle"`
		LocationID    string `json:"location_id"`
	}{}

	if err := c.BindJSON(&scanRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, Errorf("failed to parse request: %s", err))
		return
	}
	if scanRequest.StudentHandle == "" {
		c.JSON(http.StatusUnprocessableEntity, Errorf("no student handle specified"))
	}
	if scanRequest.LocationID == "" {
		c.JSON(http.StatusUnprocessableEntity, Errorf("no location specified"))
	}

	log := logrus.WithFields(logrus.Fields{
		"StudentHandle": scanRequest.StudentHandle, "LocationID": scanRequest.LocationID,
	})

	event, userError, err := trace.HandleScan(scanRequest.LocationID, scanRequest.StudentHandle)
	if err != nil {
		log.Errorf("Internal error handling scan: %s", err)
		c.JSON(http.StatusInternalServerError, Error(fmt.Errorf("internal server error: %s", err)))
		return
	}
	if userError != nil {
		log.Warnf("User error handling scan: %s", userError)
		c.JSON(http.StatusUnprocessableEntity, Error(userError))
		return
	}

	c.JSON(http.StatusCreated, Success(event))
	return
}
