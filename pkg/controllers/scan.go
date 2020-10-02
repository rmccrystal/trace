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
		LocationID    string `json:"location_id"`
	}{}

	BindJSON(c, &scanRequest)
	if scanRequest.StudentHandle == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no student handle specified")
		return
	}
	if scanRequest.LocationID == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no location specified")
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

	// Get the location and the student name
	location, found := database.DB.GetLocationByID(event.LocationID)
	if !found {
		log.Error("Could not find LocationID %s referenced by event ID %s", event.LocationID, event.ID)
		Errorf(c, http.StatusInternalServerError, "Could not find LocationID %s referenced by event ID %s", event.LocationID, event.ID)
		return
	}

	student, found := database.DB.GetStudentByID(event.StudentID)
	if !found {
		log.Error("Could not find StudentID %s referenced by event ID %s", event.StudentID, event.ID)
		Errorf(c, http.StatusInternalServerError, "Could not find StudentID %s referenced by event ID %s", event.StudentID, event.ID)
		return
	}

	// Append the location_name and the student_name
	Success(c, http.StatusCreated, struct{
		database.Event
		LocationName string `json:"location_name"`
		StudentName string `json:"student_name"`
	}{
		Event:        event,
		LocationName: location.Name,
		StudentName:  student.Name,
	})
}
