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
	visits := make([]LocationVisit, 0)
	events := database.DB.GetAllEventsBetween(minTime, maxTime)

	// create and populate latest leave and enter event
	latestLeaveEvent := make(map[database.StudentRef]database.Event)
	latestEnterEvent := make(map[database.StudentRef]database.Event)
	for _, event := range events {
		if event.EventType == database.EventLeave {
			latestLeaveEvent[event.Student] = event
		} else if event.EventType == database.EventEnter {
			latestEnterEvent[event.Student] = event
		} else {
			log.WithFields(log.Fields{
				"event": event,
			}).Errorln("invalid event type getting location visitors")
		}
	}

	for student, leaveEvent := range latestLeaveEvent {
		enterEvent, found := latestEnterEvent[student]
		if !found {
			// if there is no corresponding enter event continue
			continue
		}
		// student entered later than enter event (probably still in location)
		if leaveEvent.Time.Before(enterEvent.Time) {
			continue
		}

		visits = append(visits, LocationVisit{
			Student:   student,
			LeaveTime: leaveEvent.Time,
			EnterTime: enterEvent.Time,
		})
	}

	return visits
}
