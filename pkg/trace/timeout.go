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
	latestLeaveEvents := make(map[database.StudentRef]database.Event)
	latestEnterEvents := make(map[database.StudentRef]database.Event)

	// we have to get locations by id from the database a lot so instead
	// i'm making a cache of locations
	_locations := database.DB.GetLocations()
	locations := make(map[database.LocationRef]database.Location)
	for _, location := range _locations {
		locations[location.Ref()] = location
	}

	for _, event := range events {
		if event.EventType == database.EventLeave {
			latestLeaveEvents[event.Student] = event

			// this code doesn't really work atm for some reason but we don't really need it...
			// basically what it does is parse earlier events and add implicit logout events
			// if a student was in a location for too long but as long as the rest of the code
			// is kept running it should do the same thing

			/*
			// basically what we're doing here is checking what the current
			// latest enter event is once we hit a leave event and if the event
			// is more than the time of the leave event + the location timeout,
			// create a new leave event
			enterEvent, ok := latestEnterEvents[event.Student]
			if !ok {
				continue
			}

			location, ok := locations[enterEvent.Location]
			if !ok {
				log.WithField("event", enterEvent).Errorf("could not find location adding timeout events")
				continue
			}

			if enterEvent.Time.Add(location.Timeout).Before(event.Time) {
				newEvent := database.Event{
					Location:  enterEvent.Location,
					Student:   enterEvent.Student,
					Time:      enterEvent.Time.Add(location.Timeout - 1),
					EventType: database.EventLeave,
					Source:    database.EventSourceAutoLeave,
				}
				database.DB.CreateEvent(&newEvent)

				log.WithFields(log.Fields{
					"sourceEvent": enterEvent,
					"newEvent": newEvent,
				}).Debugln("created implicit leave event")
			}
			 */
		} else if event.EventType == database.EventEnter {
			latestEnterEvents[event.Student] = event
		} else {
			log.WithFields(log.Fields{
				"event": event,
			}).Errorln("invalid event type adding timeouts")
		}
	}

	// add the final implicit logout events to the present time
	for student, enterEvent := range latestEnterEvents {
		location, ok := locations[enterEvent.Location]
		if !ok {
			log.WithField("event", enterEvent).Errorf("could not find location adding timeout events")
			continue
		}

		// if there is an exit event after this enter event we can continue
		if leaveEvent, ok := latestLeaveEvents[student]; ok {
			if leaveEvent.Time.After(enterEvent.Time) {
				continue
			}
		}

		if enterEvent.Time.Add(location.Timeout).Before(currentTime) {
			newEvent := database.Event{
				Location:  enterEvent.Location,
				Student:   enterEvent.Student,
				Time:      enterEvent.Time.Add(location.Timeout - 1),
				EventType: database.EventLeave,
				Source:    database.EventSourceAutoLeave,
			}
			database.DB.CreateEvent(&newEvent)

			log.WithFields(log.Fields{
				"sourceEvent": enterEvent,
				"newEvent": newEvent,
			}).Debugln("created implicit leave event")
		}
	}
}

// TimeoutEventThread should be ran whenever trace is ran... it basically creates
// timeout events whenever a student times out of a location. run this on a new goroutine
// using `go TimeoutEventThread()`
func TimeoutEventThread() {
	log.Debugf("TimeoutEventThread started")
	for {
		now := time.Now()
		AddTimeoutEvents(now.Add(-6 * time.Hour), now)

		// run every 30 seconds
		time.Sleep(3 * time.Second)
	}
}
