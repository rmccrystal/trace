package trace

import (
	"fmt"
	"time"
	"trace/pkg/database"
)

func IsStudentAtLocation(studentRef database.StudentRef, locationRef database.LocationRef, t time.Time) (bool, database.Event) {
	location := locationRef.Get()

	// Get the event between the time and time - the location timeout
	lastEvent, found := database.DB.GetMostRecentEventBetween(studentRef, t.Add(location.Timeout * -1), t)

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
func GetStudentsAtLocation(locationRef database.LocationRef, t time.Time) ([]database.Student, []database.Event) {
	// iterate through all students and check if each one is at the location

	studentsAtLocation := make([]database.Student, 0)
	events := make([]database.Event, 0)

	location := locationRef.Get()

	// all events in the time frame sorted from earliest to latest
	allEvents := database.DB.GetAllEventsBetween(t.Add(location.Timeout * -2 - 1 * time.Hour), t)

	// the latest event for each student
	studentEvents := make(map[database.StudentRef]database.Event)
	for _, event := range allEvents {
		studentEvents[event.Student] = event
	}

	for student, event := range studentEvents {
		if event.EventType == database.EventEnter {
			studentsAtLocation = append(studentsAtLocation, student.Get())
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
	found = true
	return
}
