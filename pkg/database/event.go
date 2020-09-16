package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	result, err := db.Collections.Events.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}

	event.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetEvents returns a list of all events stored in the database.
// It will return an ordered list with the most recent first.
func (db *Database) GetEvents() ([]Event, error) {
	cur, err := db.Collections.Events.Find(context.TODO(), bson.D{}, &options.FindOptions{
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

// GetMostRecentEvent gets the most recent event created by the specified studentID
// If the event is found, it will be returned, otherwise, found will be false.
// If there is an error getting the most recent event, it will be returned
func (db *Database) GetMostRecentEvent(studentID primitive.ObjectID) (event Event, found bool, error error) {
	result := db.Collections.Events.FindOne(context.TODO(), bson.D{{"studentid", studentID}}, &options.FindOneOptions{
		Sort: bson.D{{"time", -1}},
	})

	err := result.Err()
	// If the event was not found
	if err == mongo.ErrNoDocuments {
		return Event{}, false, nil
	} else if err != nil { // If there was another error
		return Event{}, false, fmt.Errorf("could not get most recent event: %s", err)
	}

	if err := result.Decode(&event); err != nil {
		return Event{}, false, fmt.Errorf("error decoding event: %s", err)
	}

	return event, true, nil
}
