package trace

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"trace/pkg/database"
)

// HandleScan should be called whenever a student scans in or scans out.
// It will return the Events that it creates or an error.
// If the studentID cannot be found in the database, it will not be stored
// If there is an error with the input (locationID or studentHandle is invalid), the
// error will be returned as a userError. If there is an error accessing the database
// or any other unexpected error, it will be returned in err
func HandleScan(locationID string, studentHandle string) (ev database.Event, userError error, err error) {
	locationObjectID, _ := primitive.ObjectIDFromHex(locationID)
	location, found, err := database.DB.GetLocationByID(locationObjectID)
	if err != nil {
		return database.Event{}, nil, err
	}
	if !found {
		return database.Event{}, fmt.Errorf("location with ID %s was not found", locationID), nil
	}

	student, found, err := database.DB.GetStudentByHandle(studentHandle)
	if err != nil {
		return database.Event{}, nil, err
	}
	if !found {
		return database.Event{}, fmt.Errorf("student with handle %s was not found", studentHandle), nil
	}

	studentAtLocation, _, err := IsStudentAtLocation(student.ID, location.ID, time.Now())
	if err != nil {
		return database.Event{}, nil, fmt.Errorf("encountered error checking if student %s is at location %s", student.Name, location.Name)
	}

	// TODO: Create an implicit logout event when someone logs in again to a new location

	var eventType database.EventType
	// If the student is in the location, they are leaving, otherwise they are entering
	if studentAtLocation {
		eventType = database.EventLeave
	} else {
		eventType = database.EventEnter
	}

	event := database.Event{
		LocationID: location.ID,
		StudentID:  student.ID,
		Time:       time.Now(),
		EventType:  eventType,
		Source:     database.EventSourceScan,
	}

	if err := database.DB.CreateEvent(&event); err != nil {
		return database.Event{}, nil, fmt.Errorf("error creating event: %s", err)
	}

	// Log the event
	var evName string
	if eventType == database.EventEnter {
		evName = "into"
	} else if eventType == database.EventLeave {
		evName = "out of"
	}
	logrus.WithFields(logrus.Fields{
		"studentName": student.Name, "locationName": location.Name,
	}).Debugf("Student scanned %s the library", evName)

	return event, nil, nil
}
