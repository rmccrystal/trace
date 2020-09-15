package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EventType int32

const (
	EventEnter = iota
	EventLeave
)

// An Event represents a student either entering or leaving a location
type Event struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	LocationID primitive.ObjectID `json:"location_id"`
	StudentID  primitive.ObjectID `json:"student_id"`
	Time       time.Time          `json:"time"`
	EventType  EventType          `json:"event_type"`
}

// GetEvents gets all of the events in the Database
func (db *Database) GetEvents() []Event {
	return nil
}

// NewEvent creates an event and saves it to the Database
func (db *Database) NewEvent(location Location) *Event {
	return nil
}
