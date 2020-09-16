package database

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var TestDatabase *Database

func TestConnect(t *testing.T) {
	var err error

	TestDatabase, err = Connect(Config{
		MongoURI:     "mongodb://localhost",
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
		Time: time.Now(),
		StudentID: studentID,
	}
	event2 := Event{
		Time: time.Now().Add(-5 * time.Second),
		StudentID: studentID,
	}

	if err := TestDatabase.CreateEvent(&event1); err != nil {
		t.Fatalf("Error creating event: %s", err)
	}
	if err := TestDatabase.CreateEvent(&event2); err != nil {
		t.Fatalf("Error creating event: %s", err)
	}

	// Get the most recent event with an empty student ID
	mostRecentEvent, ok, err := TestDatabase.GetMostRecentEvent(studentID)
	if err != nil {
		t.Fatalf("Error getting most recent event: %s", err)
	}
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
	if err := TestDatabase.CreateStudent(&student); err != nil {
		t.Fatalf("Error adding student to the database: %s", err)
	}

	handle := "12345"
	foundStudent, found, err := TestDatabase.GetStudentByHandle(handle)
	if err != nil {
		t.Fatalf("Error getting student by handle: %s", err)
	}
	if !found {
		t.Fatalf("Did not find any students by the handle %s", handle)
	}
	if foundStudent.ID != student.ID {
		t.Fatalf("Found the wrong student by handle. ID should be %s, ID %s", student.ID, foundStudent.ID)
	}

	logrus.Infof("Found student %v+ using handle %s", foundStudent, handle)
}
