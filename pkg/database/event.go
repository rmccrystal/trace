package database

import (
	"context"
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
	EventSourceLoggedOut		// When a student is manually logged out through the console
	EventSourceLoggedOutAll		// When the log out all button is clicked
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
func (db *Database) CreateEvent(event *Event) {
	result, err := db.Collections.Events.InsertOne(context.TODO(), event)
	if err != nil {
		panic(err)
	}

	event.ID = result.InsertedID.(primitive.ObjectID)
}

// GetEvents returns a list of all events stored in the database.
// It will return an ordered list with the most recent first.
func (db *Database) GetEvents() []Event {
	cur, err := db.Collections.Events.Find(context.TODO(), bson.D{}, &options.FindOptions{
		// Sort by date
		Sort: bson.D{{"time", -1}},
	})
	if err != nil {
		panic(err)
	}

	events := make([]Event, 0)
	if err := cur.All(context.TODO(), &events); err != nil {
		panic(err)
	}

	return events
}

// GetMostRecentEvent gets the most recent event created by the specified studentID
// If the event is found, it will be returned, otherwise, found will be false.
// If there is an error getting the most recent event, it will be returned
func (db *Database) GetMostRecentEvent(studentID primitive.ObjectID) (event Event, found bool) {
	return db.GetMostRecentEventBetween(studentID, time.Unix(0, 0), time.Now())
}

// GetMostRecentEventBetween gets the most recent event between two time intervals
func (db *Database) GetMostRecentEventBetween(studentID primitive.ObjectID, minTime time.Time, maxTime time.Time) (event Event, found bool) {
	result := db.Collections.Events.FindOne(context.TODO(), bson.D{
		{"studentid", studentID},
		{"time", bson.M{"$lt": maxTime}},
		{"time", bson.M{"$gt": minTime}},
	}, &options.FindOneOptions{
		Sort: bson.D{{"time", -1}},
	})

	err := result.Err()
	// If the event was not found
	if err == mongo.ErrNoDocuments {
		return Event{}, false
	} else if err != nil { // If there was another error
		panic(err)
	}

	if err := result.Decode(&event); err != nil {
		panic(err)
	}

	return event, true
}

// GetAllEventsBetween gets all of the events between minTime and maxTime.
// The events will be sorted by earliest to latest.
func (db *Database) GetAllEventsBetween(minTime time.Time, maxTime time.Time) []Event {
	cursor, err := db.Collections.Events.Find(context.TODO(), bson.D{
		{"time", bson.M{"$lt": maxTime}},
		{"time", bson.M{"$gt": minTime}},
	}, &options.FindOptions{
		Sort: bson.D{{"time", 1}},
	})
	if err != nil {
		panic(err)
	}

	events := make([]Event, 0)
	if err := cursor.All(nil, &events); err != nil {
		panic(err)
	}

	return events
}