package trace

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
	"trace/pkg/database"
)

var TestDatabase *database.Database

var TestStudent *database.Student
var TestLocation *database.Location

func init() {
	logrus.SetLevel(logrus.TraceLevel)

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

	// Add an example student for the test
	TestStudent = &database.Student{
		Name:           "Ben Aaron",
		Email:          "baaron@gmail.com",
		StudentHandles: []string{"12345", "testid1"},
	}
	if err := TestDatabase.CreateStudent(TestStudent); err != nil {
		logrus.Fatalf("Error adding student to the database: %s", err)
		TestDatabase = nil
	}

	TestLocation = &database.Location{
		Name: "Library",
		Timeout: 1 * time.Hour,
	}
	if err := TestDatabase.CreateLocation(TestLocation); err != nil {
		logrus.Fatalf("Error creating location: %s", err)
		TestDatabase = nil
	}
}

func TestHandleScan(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("database not initialized")
	}

	// Test a student scanning into a location
	event, userError, err := HandleScan(TestLocation.ID.Hex(), TestStudent.StudentHandles[0])
	if err != nil {
		t.Fatalf("Error handling scan: %s", err)
	}
	if userError != nil {
		t.Fatalf("Usererror while handling valid scan: %s", userError)
	}
	if event.EventType != database.EventEnter {
		t.Fatalf("Student %s did not create an EventEnter scanning into a location", TestStudent.Name)
	}
	logrus.Infof("Successfully created event %v+ while student scanned into %s", event, TestLocation.Name)

	// Test the same student scanning out of a location
	event, userError, err = HandleScan(TestLocation.ID.Hex(), TestStudent.StudentHandles[0])
	if err != nil {
		t.Fatalf("Error handling scan: %s", err)
	}
	if userError != nil {
		t.Fatalf("Usererror while handling valid scan: %s", userError)
	}
	if event.EventType != database.EventLeave {
		t.Fatalf("Student %s did not create an EventLeave scanning out of a location", TestStudent.Name)
	}
	logrus.Infof("Successfully created event %v+ while student scanned out of %s", event, TestLocation.Name)
}

func TestIsStudentAtLocation(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("database not initialized")
	}

	// Create a test event with the student entering a location
	enterEvent := database.Event{
		LocationID: TestLocation.ID,
		StudentID:  TestStudent.ID,
		Time:       time.Now(),
		EventType:  database.EventEnter,
	}
	if err := TestDatabase.CreateEvent(&enterEvent); err != nil {
		t.Fatalf("Error creating enterEvent: %s", err)
	}

	studentAtLocation, err := IsStudentAtLocation(TestStudent.ID, TestLocation.ID, time.Now())
	if err != nil {
		t.Fatalf("Error checking if student is at location: %s", err)
	}
	if studentAtLocation != true {
		t.Fatalf("IsStudentAtLocation was false when it should be true")
	}
	logrus.Printf("Student %s was at location %s after event %v+", TestStudent.Name, TestLocation.Name, enterEvent)

	// Create a test enterEvent with the student entering a location
	leaveEvent := database.Event{
		LocationID: TestLocation.ID,
		StudentID:  TestStudent.ID,
		Time:       time.Now(),
		EventType:  database.EventLeave,
	}
	if err := TestDatabase.CreateEvent(&leaveEvent); err != nil {
		t.Fatalf("Error creating leaveEvent: %s", err)
	}

	studentAtLocation, err = IsStudentAtLocation(TestStudent.ID, TestLocation.ID, time.Now())
	if err != nil {
		t.Fatalf("Error checking if student is at location: %s", err)
	}
	if studentAtLocation != false {
		t.Fatalf("IsStudentAtLocation was true when it should be false")
	}
	logrus.Printf("Student %s was not at location %s after event %v+", TestStudent.Name, TestLocation.Name, leaveEvent)
}

func TestIsStudentAtLocationAtTime(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("database not initialized")
	}

	// Create a test event with the student entering a location
	enterEvent := database.Event{
		LocationID: TestLocation.ID,
		StudentID:  TestStudent.ID,
		Time:       time.Now(),
		EventType:  database.EventEnter,
	}
	if err := TestDatabase.CreateEvent(&enterEvent); err != nil {
		t.Fatalf("Error creating enterEvent: %s", err)
	}
	logrus.Debugf("Created enter event: %v+", enterEvent)

	// Check if students were at a location an hour ago
	studentAtLocation, err := IsStudentAtLocation(TestStudent.ID, TestLocation.ID, time.Now().Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("Error checking if student is at location: %s", err)
	}
	if studentAtLocation != false {
		t.Fatalf("IsStudentAtLocation returned true for student %s at location %s despite it checking an hour ago", TestStudent.Name, TestLocation.Name)
	}
	logrus.Infof("There were no students at location %s one hour ago", TestLocation.Name)
}

func TestGetStudentsAtLocation(t *testing.T) {
	if TestDatabase == nil {
		t.Skip("database not initialized")
	}

	// Clear the events database
	_ = database.DB.Collections.Events.Drop(nil)

	// Create a test event
	enterEvent := database.Event{
		LocationID: TestLocation.ID,
		StudentID:  TestStudent.ID,
		Time:       time.Now(),
		EventType:  database.EventEnter,
	}
	if err := TestDatabase.CreateEvent(&enterEvent); err != nil {
		t.Fatalf("Error creating enterEvent: %s", err)
	}
	logrus.Debugf("Student %s entered %s", TestStudent.Name, TestLocation.Name)

	studentsAtLocation, err := GetStudentsAtLocation(TestLocation.ID, time.Now())
	if err != nil {
		t.Fatalf("Error getting students at location: %s", err)
	}
	if studentsAtLocation[0].ID != TestStudent.ID {
		t.Fatalf("Did not corretly get the students at location %s", TestLocation.Name)
	}
	logrus.Infof("Found list of students at location %s: %v+", TestLocation.Name, studentsAtLocation)

	// Check if there are students at the location 5 hours ago. There should be none
	studentsAtLocation, err = GetStudentsAtLocation(TestLocation.ID, time.Now().Add(-5 * time.Hour))
	if err != nil {
		t.Fatalf("Error getting students at location: %s", err)
	}
	if len(studentsAtLocation) > 0 {
		t.Fatalf("Found a student at location %s 5 hours ago when the enter event was created just now", TestLocation.Name)
	}
	logrus.Infof("Found no students at location %s 5 hours ago", TestLocation.Name)
}