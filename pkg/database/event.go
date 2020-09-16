package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	LocationID primitive.ObjectID `json:"location_id"`
	StudentID  primitive.ObjectID `json:"student_id"`
	Time       time.Time          `json:"time"`
	EventType  EventType          `json:"event_type"`
	Source     EventSource		  `json:"source"`
}
