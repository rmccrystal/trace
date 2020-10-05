package trace

import (
	"fmt"
	"time"
	"trace/pkg/database"
)

// IsStudentAtLocation returns true and the corresponding event if a student is at a location at a specific time
func IsStudentAtLocation(studentRef database.StudentRef, locationRef database.LocationRef, time time.Time) (bool, database.Event) {
	location := locationRef.Get()

	// Get the event between the time and time - the location timeout
	lastEvent, found := database.DB.GetMostRecentEventBetween(studentRef, time.Add(location.Timeout * -1), time)

	if found && lastEvent.Location == locationRef {
		switch lastEvent.EventType {
		case database.EventLeave:
			return false, database.Event{}
		case database.EventEnter:
			return true, lastEvent
		default:
			panic(fmt.Sprintf("invalid event type %d", lastEvent.EventType))
		}
	}

	// If there are no past events, assume the student is not at the location
	return false, database.Event{}
}

// GetStudentsAtLocation returns a list of all students at a location at a specific time and the corresponding events.
// For most cases, the time should just be time.Now()
func GetStudentsAtLocation(locationRef database.LocationRef, time time.Time) ([]database.Student, []database.Event) {
	// iterate through all students and check if each one is at the location

	studentsAtLocation := make([]database.Student, 0)
	events := make([]database.Event, 0)

	students := database.DB.GetStudents()

	for _, student := range students {
		atLocation, event := IsStudentAtLocation(database.StudentRef(student.ID), locationRef, time)
		if atLocation {
			studentsAtLocation = append(studentsAtLocation, student)
			events = append(events, event)
		}
	}

	return studentsAtLocation, events
}

// GetStudentLocation returns the location a student is at. If the student is not at any location,
// found will be false.
func GetStudentLocation(studentRef database.StudentRef, time time.Time) (location database.Location, found bool) {
	lastEvent, found := database.DB.GetMostRecentEvent(studentRef)
	if !found {
		// If there is no most recent event for this student, we can assume they are not at a location
		return database.Location{}, false
	}

	location = lastEvent.Location.Get()

	// check lastEvent time
	if lastEvent.Time.Add(location.Timeout).Before(time) {
		return database.Location{}, false
	}

	found = true
	return
}
