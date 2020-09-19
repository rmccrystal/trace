package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"trace/pkg/database"
	"trace/pkg/trace"
)

func GetStudents(c *gin.Context) {
	locations, err := database.DB.GetStudents()
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal server error getting students: %s", locations)
		return
	}

	Success(c, http.StatusOK, locations)
}

func CreateStudent(c *gin.Context) {
	var student database.Student
	err := c.BindJSON(&student)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "failed to parse request body: %s", err)
		return
	}

	if student.Name == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no student name specified")
		return
	}

	err = database.DB.CreateStudent(&student)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal server error creating student: %s", err)
		return
	}

	Success(c, http.StatusCreated, student)
}

func GetStudentByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}
	student, found, err := database.DB.GetStudentByID(id)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "error querying database: %s", err)
		return
	}
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "student with id %s not found", id.Hex())
		return
	}

	Success(c, http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}

	success, err := database.DB.DeleteStudent(id)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal error deleting student: %s", err)
		return
	}
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find student with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, nil)
}

func UpdateStudent(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}

	var student database.Student
	err = c.BindJSON(&student)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "failed to parse request body: %s", err)
		return
	}

	success, err := database.DB.UpdateStudent(id, &student)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal error updating student: %s", err)
		return
	}
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find student with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, student)
}

func GetStudentLocation(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}
	student, found, err := database.DB.GetStudentByID(id)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "error querying database: %s", err)
		return
	}
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "student with id %s not found", id.Hex())
		return
	}

	// use current time by default
	json := struct {
		Time time.Time `json:"time"`
	}{time.Now()}
	_ = c.BindJSON(&json)

	location, found, err := trace.GetStudentLocation(student.ID, json.Time)
	if err != nil {
		Errorf(c, http.StatusInternalServerError, "internal error getting student location: %s", err)
		return
	}
	if !found {
		Success(c, http.StatusOK, nil)
		return
	}

	Success(c, http.StatusOK, location)
}