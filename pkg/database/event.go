package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// EventType determines the type of an event
type EventType int32

const (
	EventEnter = iota // When a student enters a location
	EventLeave        // When a student leaves a location
)

// EventSource is where an event came from
type EventSource int32

const (
	EventSourceScan      = iota // When a student scans in our out of a location
	EventSourceAutoLeave        // When a student leaves the library by not singing out for a period of time
)

// An Event represents a student either entering or leaving a location
type Event struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LocationID primitive.ObjectID `json:"location_id"`
	StudentID  primitive.ObjectID `json:"student_id"`
	Time       time.Time          `json:"time"`
	EventType  EventType          `json:"event_type"`
	Source     EventSource        `json:"source"`
}

// CreateEvent creates an event and adds it to the database. The
// ID element of the Event will be set if it is successful
func (db *Database) CreateEvent(event *Event) error {
	result, err := db.collections.Events.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}

	event.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetEvents returns a list of all events stored in the database.
// It will return an ordered list with the most recent first.
func (db *Database) GetEvents() ([]Event, error) {
	cur, err := db.collections.Events.Find(context.TODO(), bson.D{}, &options.FindOptions{
		// Sort by date
		Sort: bson.D{{"time", -1}},
	})
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := cur.All(context.TODO(), &events); err != nil {
		return nil, err
	}

	return events, nil
}