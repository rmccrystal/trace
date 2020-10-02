package trace

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"trace/pkg/database"
)

// IsStudentAtLocation returns true and the corresponding event if a student is at a location at a specific time
func IsStudentAtLocation(studentID primitive.ObjectID, locationID primitive.ObjectID, time time.Time) (bool, database.Event, error) {
	location, found := database.DB.GetLocationByID(locationID)
	if !found {
		return false, database.Event{}, fmt.Errorf("could not find location")
	}

	// Get the event between the time and time - the location timeout
	lastEvent, found := database.DB.GetMostRecentEventBetween(studentID, time.Add(location.Timeout * -1), time)

	if found && lastEvent.LocationID == locationID {
		switch lastEvent.EventType {
		case database.EventLeave:
			return false, database.Event{}, nil
		case database.EventEnter:
			return true, lastEvent, nil
		default:
			return false, lastEvent, fmt.Errorf("invalid event type %d", lastEvent.EventType)
		}
	}

	// If there are no past events, assume the student is not at the location
	return false, database.Event{}, nil
}

// GetStudentsAtLocation returns a list of all students at a location at a specific time and the corresponding events.
// For most cases, the time should just be time.Now()
func GetStudentsAtLocation(locationID primitive.ObjectID, time time.Time) ([]database.Student, []database.Event, error) {
	// iterate through all students and check if each one is at the location

	studentsAtLocation := make([]database.Student, 0)
	events := make([]database.Event, 0)

	location, found := database.DB.GetLocationByID(locationID)
	if !found {
		return nil, nil, fmt.Errorf("could not find location with id %s", locationID)
	}

	students := database.DB.GetStudents()

	for _, student := range students {
		atLocation, event, err := IsStudentAtLocation(student.ID, locationID, time)
		if err != nil {
			logrus.Errorf("Error checking if student %s is at location %s: %s", student.Name, location.Name, err)
		}
		if atLocation {
			studentsAtLocation = append(studentsAtLocation, student)
			events = append(events, event)
		}
	}

	return studentsAtLocation, events, nil
}

// GetStudentLocation returns the location a student is at. If the student is not at any location,
// found will be false.
func GetStudentLocation(studentID primitive.ObjectID, time time.Time) (location database.Location, found bool, err error) {
	student, found := database.DB.GetStudentByID(studentID)
	if !found {
		return database.Location{}, false, fmt.Errorf("could not find student with id %s", studentID.Hex())
	}

	lastEvent, found := database.DB.GetMostRecentEvent(studentID)
	if !found {
		// If there is no most recent event for this student, we can assume they are not at a location
		return database.Location{}, false, nil
	}

	location, found = database.DB.GetLocationByID(lastEvent.ID)
	if !found {
		return database.Location{}, false, fmt.Errorf("could not find location refrenced in student id %s", student.ID)
	}

	// check lastEvent time
	if lastEvent.Time.Add(location.Timeout).Before(time) {
		return database.Location{}, false, nil
	}

	found = true
	return
}