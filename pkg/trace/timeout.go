package trace

import (
	log "github.com/sirupsen/logrus"
	"time"
	"trace/pkg/database"
)

// AddTimeoutEvents creates leave events for enter events that have timed out
// using location.Timeout
func AddTimeoutEvents(startTime time.Time, currentTime time.Time) {
	events := database.DB.GetAllEventsBetween(startTime, currentTime)

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
			}).Errorln("invalid event type adding timeouts")
		}
	}

	for student, enterEvent := range latestEnterEvent {
		if leaveEvent, found := latestLeaveEvent[student]; found {
			// student has already left
			if leaveEvent.Time.After(enterEvent.Time) {
				continue
			}
		}

		// if the timeout time is before the current time, we should add an auto leave event
		if enterEvent.Time.Add(enterEvent.Location.Get().Timeout).Before(currentTime) {

		}
	}
}