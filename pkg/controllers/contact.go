package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"trace/pkg/database"
	"trace/pkg/trace"
)

type contact struct {
	Student         database.Student `json:"student"`
	SecondsTogether int              `json:"seconds_together"`
}

type contactReport struct {
	TargetStudent database.Student `json:"target_student"`
	StartDate     int64            `json:"start_date"`
	EndDate       int64            `json:"end_date"`
	Contacts      []contact        `json:"contacts"`
}

// GET /api/v1/trace/:id
func GenerateContactReport(c *gin.Context) {
	student, err := database.DB.GetStudentByIDString(c.Param("id"))
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	scanRequest := struct {
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
	}{}

	if !BindJSON(c, &scanRequest) {
		return
	}

	report, err := trace.GenerateContactReport(&student, time.Unix(scanRequest.StartTime, 0), time.Unix(scanRequest.EndTime, 0), 1)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	newReport := contactReport{TargetStudent: student, StartDate: scanRequest.StartTime, EndDate: scanRequest.EndTime}
	for s, t := range report.Contacts[0] {
		newReport.Contacts = append(newReport.Contacts, contact{Student: s.Get(), SecondsTogether: int(t.Seconds())})
	}

	Success(c, http.StatusOK, newReport)
}
