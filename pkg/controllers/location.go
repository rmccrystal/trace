package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"trace/pkg/database"
	"trace/pkg/trace"
)

func GetLocations(c *gin.Context) {
	locations := database.DB.GetLocations()
	Success(c, http.StatusOK, locations)
}

func CreateLocation(c *gin.Context) {
	var location database.Location
	if success := BindJSON(c, &location); !success {
		return
	}

	if location.Name == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no location name specified")
		return
	}

	database.DB.CreateLocation(&location)
	Success(c, http.StatusCreated, location)
}

func GetLocationByID(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	location, found := database.DB.GetLocationByID(id)
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "location with id %s not found", id.Hex())
		return
	}

	Success(c, http.StatusOK, location)
}

func DeleteLocation(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	success = database.DB.DeleteLocation(id)
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find location with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, nil)
}

func UpdateLocation(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	var location database.Location
	if success := BindJSON(c, &location); !success {
		return
	}

	success = database.DB.UpdateLocation(id, &location)
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find location with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, location)
}

func GetStudentsAtLocation(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	// use current time by default
	json := struct {
		Time time.Time `json:"time"`
	}{time.Now()}
	_ = c.ShouldBindJSON(&json)

	students, events, err := trace.GetStudentsAtLocation(id, json.Time)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not get students at location: %s", err)
		return
	}

	/* Create a json response formatted as:
	[{
		"student": (student),
		"time": (time)
	}]
	*/
	var resp = make([]map[string]interface{}, 0)
	for i := range students {
		resp = append(resp, map[string]interface{}{
			"student": students[i],
			"time":    events[i].Time,
		})
	}

	Success(c, http.StatusOK, resp)
}

func LogoutAllStudentsAtLocation(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	location, found := database.DB.GetLocationByID(id)
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "location with id %s not found", id.Hex())
		return
	}

	students, _, err := trace.GetStudentsAtLocation(location.ID, time.Now())
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "error getting students at location: %s", err)
		return
	}

	for _, student := range students {
		newEvent := database.Event{
			Location:  database.LocationRef(location.ID),
			Student:   database.StudentRef(student.ID),
			Time:      time.Now(),
			EventType: database.EventLeave,
			Source:    database.EventSourceLoggedOutAll,
		}
		database.DB.CreateEvent(&newEvent)
	}

	Success(c, http.StatusCreated, nil)
}
