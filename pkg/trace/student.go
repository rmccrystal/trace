package trace

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"trace/pkg/database"
)

// IsStudentAtLocation returns true if a student is at a location at a specific time
func IsStudentAtLocation(studentID primitive.ObjectID, locationID primitive.ObjectID, time time.Time) (bool, error) {
	location, found, err := database.DB.GetLocationByID(locationID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, fmt.Errorf("could not find location")
	}

	// Get the event between the time and time - the location timeout
	lastEvent, found, err := database.DB.GetMostRecentEventBetween(studentID, time.Add(location.Timeout * -1), time)
	if err != nil {
		return false, err
	}

	if found && lastEvent.LocationID == locationID {
		switch lastEvent.EventType {
		case database.EventLeave:
			return false, nil
		case database.EventEnter:
			return true, nil
		default:
			return false, fmt.Errorf("invalid event type %d", lastEvent.EventType)
		}
	}

	// If there are no past events, assume the student is not at the location
	return false, nil
}

// GetStudentsAtLocation returns a list of all students at a location at a specific time.
// For most cases, the time should just be time.Now()
func GetStudentsAtLocation(locationID primitive.ObjectID, time time.Time) ([]database.Student, error) {
	// iterate through all students and check if each one is at the location

	studentsAtLocation := make([]database.Student, 0)

	location, found, err := database.DB.GetLocationByID(locationID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("could not find location with id %s", locationID)
	}

	students, err := database.DB.GetStudents()
	if err != nil {
		return nil, fmt.Errorf("error getting students: %s", err)
	}

	for _, student := range students {
		atLocation, err := IsStudentAtLocation(student.ID, locationID, time)
		if err != nil {
			logrus.Errorf("Error checking if student %s is at location %s: %s", student.Name, location.Name, err)
		}
		if atLocation {
			studentsAtLocation = append(studentsAtLocation, student)
		}
	}

	return studentsAtLocation, nil
}

// GetStudentLocation returns the location a student is at. If the student is not at any location,
// found will be false.
func GetStudentLocation(studentID primitive.ObjectID, time time.Time) (location database.Location, found bool, err error) {
	student, found, err := database.DB.GetStudentByID(studentID)
	if err != nil {
		return database.Location{}, false, err
	}
	if !found {
		return database.Location{}, false, fmt.Errorf("could not find student with id %s", studentID.Hex())
	}

	lastEvent, found, err := database.DB.GetMostRecentEvent(studentID)
	if err != nil {
		return database.Location{}, false, fmt.Errorf("error getting most recent event: %s", err)
	}
	if !found {
		// If there is no most recent event for this student, we can assume they are not at a location
		return database.Location{}, false, nil
	}

	location, found, err = database.DB.GetLocationByID(lastEvent.ID)
	if err != nil {
		return database.Location{}, false, fmt.Errorf("error getting location: %s", err)
	}
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