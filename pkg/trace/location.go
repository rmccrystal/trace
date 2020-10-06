package trace

import (
	log "github.com/sirupsen/logrus"
	"time"
	"trace/pkg/database"
)

type LocationVisit struct {
	Student   database.StudentRef `json:"student"`
	LeaveTime time.Time           `json:"leave_time"`
	EnterTime time.Time           `json:"enter_time"`
}

// GetLocationVisitors returns a list of LocationVisit objects representing
// who has entered a location in the time range. This will not return students who
// are currently in the location
func GetLocationVisitors(locationRef database.LocationRef, minTime time.Time, maxTime time.Time) []LocationVisit {
	students := database.DB.GetStudents()
	visits := make([]LocationVisit, 0)
	for _, student := range students {
		// iterate through every student and check if they left the location
		leaveEvent, found := database.DB.GetMostRecentEventBetween(student.Ref(), minTime, maxTime)
		if !found {
			continue
		}

		// check if they were leaving
		if leaveEvent.EventType != database.EventLeave {
			continue
		}

		enterEvent, found := database.DB.GetMostRecentEventBetweenWithType(student.Ref(), minTime, maxTime, database.EventEnter)
		if !found {
			log.WithField("leaveEvent", leaveEvent).Errorf("Found a logout leaveEvent without a corresponding login leaveEvent")
			enterEvent.Time = leaveEvent.Time
		}

		visits = append(visits, LocationVisit{
			Student:   student.Ref(),
			LeaveTime: leaveEvent.Time,
			EnterTime: enterEvent.Time,
		})
	}

	return visits
}
