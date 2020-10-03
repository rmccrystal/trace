package database

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// A Location represents any location at the school that can be signed in or out
type Location struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `json:"name"`
	// The time it takes for a student to automatically time out
	Timeout time.Duration      `json:"timeout"`
}

// CreateLocation creates a location and adds it to the database. The
// ID element of the Location will be set if it is successful
func (db *Database) CreateLocation(location *Location) {
	if location.Timeout == 0 {
		logrus.Warnf("Created a location with 0 timeout. Defaulting to 2 hours")
		location.Timeout = 2 * time.Hour
	}

	result, err := db.Collections.Locations.InsertOne(context.TODO(), location)
	if err != nil {
		panic(err)
	}

	location.ID = result.InsertedID.(primitive.ObjectID)
}

// GetLocationByID gets a location by its ID. If not found, found will be false and
// err will be nil.
func (db *Database) GetLocationByID(id primitive.ObjectID) (location Location, found bool) {
	result := db.Collections.Locations.FindOne(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Location{}, false
		}
		panic(err)
	}

	if err := result.Decode(&location); err != nil {
		panic(err)
	}

	found = true
	return
}

// GetLocations returns a list of all locations stored in the database.
func (db *Database) GetLocations() []Location {
	cur, err := db.Collections.Locations.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	locations := make([]Location, 0)
	if err := cur.All(context.TODO(), &locations); err != nil {
		panic(err)
	}

	return locations
}

// DeleteLocation deletes a location from the database by ID. If the location could not be
// found, success will be false
func (db *Database) DeleteLocation(id primitive.ObjectID) bool {
	result := db.Collections.Locations.FindOneAndDelete(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}

	return true
}

// UpdateLocation finds a location by its ID and updates it
func (db *Database) UpdateLocation(id primitive.ObjectID, newLocation *Location) bool {
	result := db.Collections.Locations.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{"$set": newLocation})
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}

	if err := result.Decode(newLocation); err != nil {
		panic(err)
	}

	return true
}