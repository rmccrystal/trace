package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

	// Get the mongo URI from env var
	mongoURI, found := os.LookupEnv("TEST_MONGO_URI")
	// Set it to localhost by default
	if !found {
		mongoURI = "mongodb://localhost"
	}

	TestDatabase, err = database.Connect(database.Config{
		MongoURI:     mongoURI,
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
		StudentHandles: []string{"testhandle"},
	}
	if err := TestDatabase.CreateStudent(TestStudent); err != nil {
		logrus.Fatalf("Error adding student to the database: %s", err)
		TestDatabase = nil
	}

	TestLocation = &database.Location{
		Name:    "Library",
		Timeout: 1 * time.Hour,
	}
	if err := TestDatabase.CreateLocation(TestLocation); err != nil {
		logrus.Fatalf("Error creating location: %s", err)
		TestDatabase = nil
	}
}

func sendTestRequest(handler func(ctx *gin.Context), json []byte) (code int, body []byte) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// The method and url don't matter because we're running the handler function directly
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	handler(c)

	body, _ = ioutil.ReadAll(w.Body)

	return w.Code, body
}

func TestOnScan(t *testing.T) {
	code, resp := sendTestRequest(OnScan, []byte(fmt.Sprintf(`
{
	"student_handle": "testhandle",
	"location_id": "%s"
}`, TestLocation.ID.Hex())))

	assert.Equal(t, 201, code)

	var createdEvent database.Event
	assert.NoError(t, json.Unmarshal(resp, &createdEvent), "failed to unmarshal response")

	mostRecentEvent, found, err := TestDatabase.GetMostRecentEvent(TestStudent.ID)
	assert.NoError(t, err, "error getting most recent event")
	assert.Truef(t, found, "no event was created")
	assert.Equalf(t, mostRecentEvent.ID, createdEvent.ID, "the returned event id %s did not match the most recent event id %s", createdEvent.ID, mostRecentEvent.ID)
	assert.Equal(t, createdEvent.EventType, database.EventEnter, "incorrect event type %s was created", createdEvent.EventType)
}
