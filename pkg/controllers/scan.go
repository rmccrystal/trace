package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"trace/pkg/database"
	"trace/pkg/trace"
)

// POST /api/v1/scan
// Called whenever someone scans their barcode
func OnScan(c *gin.Context) {
	scanRequest := struct {
		StudentHandle string `json:"student_handle"`
		LocationID    database.LocationRef `json:"location_id"`
	}{}

	BindJSON(c, &scanRequest)
	if scanRequest.StudentHandle == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no student handle specified")
		return
	}

	log := logrus.WithFields(logrus.Fields{
		"StudentHandle": scanRequest.StudentHandle, "LocationID": scanRequest.LocationID,
	})

	event, userError, err := trace.HandleScan(scanRequest.LocationID, scanRequest.StudentHandle)
	if err != nil {
		log.Errorf("Internal error handling scan: %s", err)
		Errorf(c, http.StatusInternalServerError, "internal server error: %s", err)
		return
	}
	if userError != nil {
		log.Warnf("User error handling scan: %s", userError)
		Errorf(c, http.StatusUnprocessableEntity, "%s", userError)
		return
	}

	Success(c, http.StatusCreated, event)
}
