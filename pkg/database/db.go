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

	collections struct {
		Event    *mongo.Collection
		Location *mongo.Collection
		Student  *mongo.Collection
	}
}

// Connect connects to the Database using a DatabaseConfig
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

	database := client.Database(config.DatabaseName)

	log.Info("Connected to Database")

	DB = &Database{Client: client, Database: database}

	return DB, nil
}
