package trace

import (
	"fmt"
	"time"
	"trace/pkg/database"
)

// HandleScan should be called whenever a student scans in or scans out.
// It will return the Events that it creates or an error.
// If the studentID cannot be found in the database, it will not be stored
func HandleScan(locationID string, studentHandle string) (database.Event, error) {
	location, found, err := database.DB.GetLocationByID(locationID)
	if err != nil {
		return database.Event{}, err
	}
	if !found {
		return database.Event{}, fmt.Errorf("location with ID %s was not found", locationID)
	}

	student, found, err := database.DB.GetStudentByHandle(studentHandle)
	if err != nil {
		return database.Event{}, err
	}
	if !found {
		return database.Event{}, fmt.Errorf("student with handle %s was not found", studentHandle)
	}

	// If the student is entering or leaving the location
	var eventType database.EventType

	lastEvent, found, _ := database.DB.GetMostRecentEvent(student.ID)

	// If there are no past events, assume the student is entering the location
	if !found {
		eventType = database.EventEnter
	} else {
		// If the student last entered the library, they are leaving
		// TODO: If the student entered the library like a day ago then they should be entering again
		switch lastEvent.EventType {
		case database.EventLeave:
			eventType = database.EventEnter
		case database.EventEnter:
			eventType = database.EventLeave
		}
	}

	event := database.Event{
		LocationID: location.ID,
		StudentID:  student.ID,
		Time:       time.Now(),
		EventType:  eventType,
		Source:     database.EventSourceScan,
	}

	if err := database.DB.CreateEvent(&event); err != nil {
		return database.Event{}, fmt.Errorf("error creating event: %s", err)
	}

	return event, nil
}
