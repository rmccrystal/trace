package trace

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
	"trace/pkg/database"
)

// HandleScan should be called whenever a student scans in or scans out.
// It will return the Events that it creates or an error.
// If the studentID cannot be found in the database, it will not be stored
// If there is an error with the input (locationID or studentHandle is invalid), the
// error will be returned as a userError. If there is an error accessing the database
// or any other unexpected error, it will be returned in err
func HandleScan(locationRef database.LocationRef, studentHandle string) (ev database.Event, userError error, err error) {
	location := locationRef.Get()

	student, found := database.DB.GetStudentByHandle(studentHandle)
	if !found {
		return database.Event{}, fmt.Errorf("student with handle %s was not found", studentHandle), nil
	}

	studentAtLocation, _ := IsStudentAtLocation(student.Ref(), location.Ref(), time.Now())

	var eventType database.EventType
	// If the student is in the location, they are leaving, otherwise they are entering
	if studentAtLocation {
		eventType = database.EventLeave
	} else {
		eventType = database.EventEnter
	}

	event := database.Event{
		Location: database.LocationRef(location.ID),
		Student:  database.StudentRef(student.ID),
		Time:       time.Now(),
		EventType:  eventType,
		Source:     database.EventSourceScan,
	}

	database.DB.CreateEvent(&event)

	// Log the event
	var evName string
	if eventType == database.EventEnter {
		evName = "into"
	} else if eventType == database.EventLeave {
		evName = "out of"
	}
	logrus.WithFields(logrus.Fields{
		"studentName": student.Name, "locationName": location.Name,
	}).Debugf("Student scanned %s a location", evName)

	return event, nil, nil
}
