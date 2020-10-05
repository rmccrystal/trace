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
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	Success(c, http.StatusOK, location)
}

func DeleteLocation(c *gin.Context) {
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	_ = database.DB.DeleteLocation(location.ID)

	Success(c, http.StatusOK, nil)
}

func UpdateLocation(c *gin.Context) {
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	var newLocation database.Location
	if success := BindJSON(c, &location); !success {
		return
	}

	_ = database.DB.UpdateLocation(location.ID, &newLocation)

	Success(c, http.StatusOK, newLocation)
}

func GetStudentsAtLocation(c *gin.Context) {
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	// use current time by default
	json := struct {
		Time time.Time `json:"time"`
	}{time.Now()}
	_ = c.ShouldBindJSON(&json)

	students, events := trace.GetStudentsAtLocation(location.Ref(), json.Time)

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
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	students, _ := trace.GetStudentsAtLocation(location.Ref(), time.Now())

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

func VisitedLocationToday(c *gin.Context) {
	location, err := database.DB.GetLocationByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	visitReport := trace.GetLocationVisitors(location.Ref(), time.Now().Add(-24 * time.Hour), time.Now())

	Success(c, http.StatusOK, visitReport)
}