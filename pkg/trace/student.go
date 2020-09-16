package trace

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"trace/pkg/database"
)

// IsStudentAtLocation returns true if a student is at a location
func IsStudentAtLocation(studentID primitive.ObjectID, locationID primitive.ObjectID) (bool, error) {
	lastEvent, found, err := database.DB.GetMostRecentEvent(studentID)
	if err != nil {
		return false, err
	}

	// If there are no past events, assume the student is not at the location
	if !found {
		return false, nil
	} else {
		// TODO: If the student entered the library like a day ago and never left this should be false
		switch lastEvent.EventType {
		case database.EventLeave:
			return false, nil
		case database.EventEnter:
			return true, nil
		default:
			return false, fmt.Errorf("invalid event type %d", lastEvent.EventType)
		}
	}
}
