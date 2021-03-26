package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"trace/pkg/api"
	"trace/pkg/database"
	"trace/pkg/trace"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	addr := envOr("LISTEN_ADDRESS", "0.0.0.0:8080")
	mongoUri := envOr("MONGO_URI", "mongodb://localhost")
	databaseName := envOr("DATABASE_NAME", "prod")

	config := api.Config{
		DatabaseConfig: database.Config{
			MongoURI:     mongoUri,
			DatabaseName: databaseName,
		},
		Timeout:        3 * time.Hour,
	}

	// we're using the global database
	_, err := database.Connect(config.DatabaseConfig)
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	// see the function docs for more info (trace/timeout.go)
	go trace.TimeoutEventThread()

	if err := api.Listen(addr, &config); err != nil {
		logrus.Fatalf("Failed to listen on %s: %s", addr, err)
	}
}

// envOr returns an env variable by its key or the defaultValue if it is not found
func envOr(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return value
}