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
	locations := database.DB.GetStudents()
	Success(c, http.StatusOK, locations)
}

func LogoutStudent(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	student, found := database.DB.GetStudentByID(id)
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "student with id %s not found", id.Hex())
		return
	}

	body := struct{
		LocationID string `json:"location_id"`
	}{}
	if success := BindJSON(c, &body); !success {
		return
	}

	// TODO: use GetIDParam or modify the func?
	locationID, err := primitive.ObjectIDFromHex(body.LocationID)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return
	}
	location, found := database.DB.GetLocationByID(locationID)
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "location with id %s not found", id.Hex())
		return
	}

	newEvent := database.Event{
		LocationID: location.ID,
		StudentID:  student.ID,
		Time:       time.Now(),
		EventType:  database.EventLeave,
		Source:     database.EventSourceLoggedOut,
	}
	database.DB.CreateEvent(&newEvent)

	Success(c, http.StatusCreated, newEvent)
}

func CreateStudent(c *gin.Context) {
	var student database.Student

	if success := BindJSON(c, &student); !success {
		return
	}

	if student.Name == "" {
		Errorf(c, http.StatusUnprocessableEntity, "no student name specified")
		return
	}

	database.DB.CreateStudent(&student)

	Success(c, http.StatusCreated, student)
}

func CreateStudents(c *gin.Context) {
	var students []database.Student
	if success := BindJSON(c, &students); !success {
		return
	}

	for i := range students {
		if students[i].Name == "" {
			Errorf(c, http.StatusUnprocessableEntity, "no student name specified")
			return
		}

		database.DB.CreateStudent(&students[i])

	}

	Success(c, http.StatusCreated, students)
}

func GetStudentByID(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	student, found := database.DB.GetStudentByID(id)
	if !found {
		Errorf(c, http.StatusUnprocessableEntity, "student with id %s not found", id.Hex())
		return
	}

	Success(c, http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	success = database.DB.DeleteStudent(id)
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find student with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, nil)
}

func UpdateStudent(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	var student database.Student
	if success := BindJSON(c, &student); !success {
		return
	}

	success = database.DB.UpdateStudent(id, &student)
	if !success {
		Errorf(c, http.StatusUnprocessableEntity, "could not find student with ID %s", id.Hex())
		return
	}

	Success(c, http.StatusOK, student)
}

func GetStudentLocation(c *gin.Context) {
	id, success := GetIDParam(c)
	if !success {
		return
	}

	student, found := database.DB.GetStudentByID(id)

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