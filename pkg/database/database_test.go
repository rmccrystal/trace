package database

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
	"time"
)

var TestDatabase *Database

func TestConnect(t *testing.T) {
	var err error

	// Get the mongo URI from env var
	mongoURI, found := os.LookupEnv("TEST_MONGO_URI")
	// Set it to localhost by default
	if !found {
		mongoURI = "mongodb://localhost"
	}

	TestDatabase, err = Connect(Config{
		MongoURI:     mongoURI,
		DatabaseName: "tests",
	})

	if err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	// Purge the test database
	_ = TestDatabase.Database.Drop(nil)
	logrus.Infof("Purged the tests database")
}

func TestDatabase_GetMostRecentEvent(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("TestDatabase was nil")
	}

	studentID, _ := primitive.ObjectIDFromHex("000000000000000000000001")
	event1 := Event{
		Time:      time.Now(),
		StudentID: studentID,
	}
	event2 := Event{
		Time:      time.Now().Add(-5 * time.Second),
		StudentID: studentID,
	}

	TestDatabase.CreateEvent(&event1)
	TestDatabase.CreateEvent(&event2)

	// Get the most recent event with an empty student ID
	mostRecentEvent, ok := TestDatabase.GetMostRecentEvent(studentID)
	if !ok {
		t.Fatalf("Could not find the most recent event")
	}

	if mostRecentEvent.ID != event1.ID {
		t.Fatalf("The mostRecentEvent ID was %s while it should have been %s", mostRecentEvent.ID, event1.ID)
	}

	logrus.Infof("Found most recent event for student with mongo ID %s: %v+", studentID, mostRecentEvent)
}

func TestDatabase_GetStudentByHandle(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("TestDatabase was nil")
	}

	// Add an example student for the test
	student := Student{
		Name:           "Ben Aaron",
		Email:          "baaron@gmail.com",
		StudentHandles: []string{"12345", "testid1"},
	}
	TestDatabase.CreateStudent(&student)

	handle := "12345"
	foundStudent, found := TestDatabase.GetStudentByHandle(handle)
	if !found {
		t.Fatalf("Did not find any students by the handle %s", handle)
	}
	if foundStudent.ID != student.ID {
		t.Fatalf("Found the wrong student by handle. ID should be %s, ID %s", student.ID, foundStudent.ID)
	}

	logrus.Infof("Found student %v+ using handle %s", foundStudent, handle)
}
