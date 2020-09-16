package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// A Location represents any location at the school that can be signed in or out
type Location struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `json:"name"`
}

// GetLocations returns all of the locations stored in the Database
func (db *Database) GetLocations() []Location {
	return nil
}
