package trace

import (
	"time"
	"trace/pkg/database"
)

type LocationVisit struct {
	Student   database.StudentRef `json:"student"`
	LeaveTime time.Time           `json:"leave_time"`
}

// GetLocationVisitors returns a list of LocationVisit objects representing
// who has entered a location in the time range. This will not return students who
// are currently in the location
func GetLocationVisitors(locationRef database.LocationRef, minTime time.Time, maxTime time.Time) []LocationVisit {
	students := database.DB.GetStudents()
	var visits []LocationVisit
	for _, student := range students {
		// iterate through every student and check if they left the location
		event, found := database.DB.GetMostRecentEventBetween(student.Ref(), minTime, maxTime)
		if !found {
			continue
		}

		// check if they were leaving
		if event.EventType != database.EventLeave {
			continue
		}

		visits = append(visits, LocationVisit{
			Student:   student.Ref(),
			LeaveTime: event.Time,
		})
	}

	return visits
}
