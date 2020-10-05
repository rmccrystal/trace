package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// A Location represents any location at the school that can be signed in or out
type Location struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `json:"name"`
	// The time it takes for a student to automatically time out
	Timeout time.Duration      `json:"timeout"`
}