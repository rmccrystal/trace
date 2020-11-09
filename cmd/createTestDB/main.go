package main

import (
	"github.com/sirupsen/logrus"
	"time"
	"trace/pkg/database"
)

func main() {
	databaseConfig := database.Config{
		MongoURI:     "mongodb://localhost",
		DatabaseName: "dev",
	}

	db, err := database.Connect(databaseConfig)
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	_ = db.Database.Drop(nil)

	students := map[string][]string {
		"Ryan McCrystal": {"ryan", "_ryan"},
		"Ben Aaron": {"ben"},
		"Cai Noel": {"cai"},
	}

	locations := []string {
		"Library", "Gym", "Test Location 1", "Front Office",
	}

	for name, handles := range students {
		newStudent := database.Student{
			Name:           name,
			StudentHandles: handles,
		}
		db.CreateStudent(&newStudent)
	}

	for _, locationName := range locations {
		db.CreateLocation(&database.Location{
			Timeout: 4 * time.Hour,
			Name: locationName,
		})
	}
}
