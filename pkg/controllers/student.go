package controllers

import (
	"github.com/gin-gonic/gin"
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
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	body := struct {
		LocationID database.LocationRef `json:"location_id"`
	}{}
	if success := BindJSON(c, &body); !success {
		return
	}

	newEvent := database.Event{
		Location:  body.LocationID,
		Student:   student.Ref(),
		Time:      time.Now(),
		EventType: database.EventLeave,
		Source:    database.EventSourceLoggedOut,
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
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	Success(c, http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	_ = database.DB.DeleteStudent(student.ID)

	Success(c, http.StatusOK, nil)
}

func UpdateStudent(c *gin.Context) {
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	var newStudent database.Student
	if success := BindJSON(c, &newStudent); !success {
		return
	}

	_ = database.DB.UpdateStudent(student.ID, &newStudent)

	Success(c, http.StatusOK, student)
}

func GetStudentLocation(c *gin.Context) {
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	// use current time by default
	json := struct {
		Time time.Time `json:"time"`
	}{time.Now()}
	_ = c.BindJSON(&json)

	location, found := trace.GetStudentLocation(database.StudentRef(student.ID), json.Time)
	if !found {
		Success(c, http.StatusOK, nil)
		return
	}

	Success(c, http.StatusOK, location)
}

func GetStudentContacts(c *gin.Context) {
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	body := struct {
		// how much time we want to go back and look for contacts
		Duration time.Duration `json:"duration"`
		MaxDepth int           `json:"max_depth"`
	}{}
	if err := c.BindJSON(&body); err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	report, err := trace.GenerateContactReport(&student, time.Now().Add(-1*body.Duration), time.Now(), body.MaxDepth)
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	// since the result contains a map that we cannot encode as json, we need to create a different struct for the response
	// doing this now i kind of wish that I used node / ts for the backend instead of this... much better interop with json
	resp := struct {
		TargetStudent database.StudentRef `json:"target_student"`
		Contacts      []struct {
			Student  database.StudentRef `json:"student"`
			Duration time.Duration       `json:"duration"`
			Depth    int                 `json:"depth"`
		} `json:"contacts"`
	}{TargetStudent: report.TargetStudent.Ref()}

	// populate the struct... this is really messy and i don't like this but it works
	for depth, contact := range report.Contacts {
		for student, duration := range contact {
			resp.Contacts = append(resp.Contacts, struct {
				Student  database.StudentRef `json:"student"`
				Duration time.Duration       `json:"duration"`
				Depth    int                 `json:"depth"`
			}{Student: student, Duration: duration, Depth: depth})
		}
	}

	Success(c, http.StatusOK, resp)
}
