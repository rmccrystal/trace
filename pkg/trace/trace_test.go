package trace

import (
	"github.com/sirupsen/logrus"
	"testing"
	"trace/pkg/database"
)

var TestDatabase *database.Database

func init()  {
	var err error

	TestDatabase, err = database.Connect(database.Config{
		MongoURI:     "mongodb://localhost",
		DatabaseName: "tests",
	})

	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	// Purge the test database
	_ = TestDatabase.Database.Drop(nil)
	logrus.Infof("Purged the tests database")
}

func TestHandleScan(t *testing.T) {
	
}
