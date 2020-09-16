// database contains functions to interact with the stored contact tracing data
package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DB is a global Database that will be set after Connect is called
var DB *Database

// DatabaseConfig configures the connection to a Database
type Config struct {
	// The connecting string to the mongo Database
	// See https://docs.mongodb.com/manual/reference/connection-string/
	MongoURI string `json:"mongo_uri"`

	// The name of the Database to be used
	DatabaseName string `json:"database_name"`
}

// A Database manages all of the models
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
	config   Config

	Collections struct {
		Events    *mongo.Collection
		Locations *mongo.Collection
		Students  *mongo.Collection
	}
}

func newDatabase(client *mongo.Client, config Config) (Database, error) {
	database := Database{Client: client, config: config}

	// Create the mongo database
	database.Database = client.Database(config.DatabaseName)

	// Create the Collections
	database.Collections.Events = database.Database.Collection("events")
	database.Collections.Locations = database.Database.Collection("locations")
	database.Collections.Students = database.Database.Collection("students")

	return database, nil
}

// Connect connects to the Database using a DatabaseConfig.
// If successful, the global DB object will be set.
func Connect(config Config) (*Database, error) {
	// Set the Client options
	clientOptions := options.Client().ApplyURI(config.MongoURI)

	// Create the context to time out after some seconds
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	// Connect to the mongodb Database
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Try pinging the Database to see if it connects
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("could not ping Database: %s", err)
	}

	log.Info("Connected to Database")

	// Create the database object
	database, err := newDatabase(client, config)
	if err != nil {
		return nil, err
	}

	// Warn if we already set the mongo
	if DB != nil {
		log.Errorf("Connected to the mongo while it was already connected. Continuing.")
	}

	DB = &database

	return &database, nil
}
