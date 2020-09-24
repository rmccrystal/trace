package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"trace/pkg/database"
	"trace/pkg/trace"
)

func GetLocations(c *gin.Context) {
	locations, err := database.DB.GetLocations()
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal server error getting locations: %s", locations)
		return
	}

	Success(c, http.StatusOK, locations)
}

func CreateLocation(c *gin.Context) {
	var location database.Location
	err := c.BindJSON(&location)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "failed to parse request body: %s", err)
		return
	}

	if location.Name == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no location name specified")
		return
	}

	err = database.DB.CreateLocation(&location)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal server error creating location: %s", err)
		return
	}

	Success(c, http.StatusCreated, location)
}

func GetLocationByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}
	location, found, err := database.DB.GetLocationByID(id)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "error querying database: %s", err)
		return
	}
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "location with id %s not found", id.Hex())
		return
	}

	Success(c, http.StatusOK, location)
}

func DeleteLocation(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}

	success, err := database.DB.DeleteLocation(id)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal error deleting location: %s", err)
		return
	}
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find location with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, nil)
}

func UpdateLocation(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}

	var location database.Location
	err = c.BindJSON(&location)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "failed to parse request body: %s", err)
		return
	}

	success, err := database.DB.UpdateLocation(id, &location)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal error updating location: %s", err)
		return
	}
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find location with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, location)
}

func GetStudentsAtLocation(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}

	// use current time by default
	json := struct {
		Time time.Time `json:"time"`
	}{time.Now()}
	_ = c.ShouldBindJSON(&json)

	students, err := trace.GetStudentsAtLocation(id, json.Time)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not get students at location: %s", err)
		return
	}

	Success(c, http.StatusOK, students)
}
