package main

import (
	"github.com/sirupsen/logrus"
	"time"
	"trace/pkg/api"
	"trace/pkg/database"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	addr := "0.0.0.0:8080"
	config := api.Config{
		DatabaseConfig: database.Config{
			MongoURI:     "mongodb://localhost",
			DatabaseName: "dev",
		},
		Timeout:        3 * time.Hour,
	}

	// we're using the global database
	_, err := database.Connect(config.DatabaseConfig)
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	if err := api.Listen(addr, &config); err != nil {
		logrus.Fatalf("Failed to listen on %s: %s", addr, err)
	}
}
