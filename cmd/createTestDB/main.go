package main

import (
	"github.com/sirupsen/logrus"
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
		"Mrs. Madden": {"madden"},
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
		if err := db.CreateStudent(&newStudent); err != nil {
			panic(err)
		}
	}

	for _, locationName := range locations {
		if err := db.CreateLocation(&database.Location{
			Name: locationName,
		}); err != nil {
			panic(err)
		}
	}
}
