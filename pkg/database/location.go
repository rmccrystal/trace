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
func (db *Database) CreateLocation(location *Location) error {
	if location.Timeout == 0 {
		logrus.Warnf("Created a location with 0 timeout. Defaulting to 2 hours")
		location.Timeout = 2 * time.Hour
	}

	result, err := db.Collections.Locations.InsertOne(context.TODO(), location)
	if err != nil {
		return err
	}

	location.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetLocationByID gets a location by its ID. If not found, found will be false and
// err will be nil. If there is an error retrieving the location, it will be returned
func (db *Database) GetLocationByID(id primitive.ObjectID) (location Location, found bool, err error) {
	result := db.Collections.Locations.FindOne(nil, bson.M{"_id": id})

	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Location{}, false, nil
		}
		return Location{}, false, err
	}

	if err := result.Decode(&location); err != nil {
		return Location{}, true, err
	}

	found = true
	return
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
