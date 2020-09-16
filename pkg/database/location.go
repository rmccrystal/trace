package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A Location represents any location at the school that can be signed in or out
type Location struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `json:"name"`
}

// CreateLocation creates a location and adds it to the database. The
// ID element of the Location will be set if it is successful
func (db *Database) CreateLocation(location *Location) error {
	result, err := db.Collections.Locations.InsertOne(context.TODO(), location)
	if err != nil {
		return err
	}

	location.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetLocations returns a list of all locations stored in the database.
func (db *Database) GetLocations() ([]Location, error) {
	cur, err := db.Collections.Events.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := cur.All(context.TODO(), &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

