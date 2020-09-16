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

	if found {
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

	var studentsAtLocation []database.Student

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