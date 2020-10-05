package api

import (
	"errors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"trace/pkg/controllers"
	"trace/pkg/database"
)

const frontendDirectory = "frontend/build"

// Listen starts the api server at the specified address
func Listen(addr string, config *Config) error {
	if database.DB == nil {
		return errors.New("database is not connected")
	}

	GlobalConfig = config

	r := gin.Default()

	r.Use(gin.Recovery())

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if username != "" && password != "" {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"password": password,
		}).Debugf("Enabling authentication")

		r.Use(gin.BasicAuth(gin.Accounts{
			username: password,
		}))
	}

	api := r.Group("/api")

	api.POST("scan", controllers.OnScan)

	api.POST("location", controllers.CreateLocation)
	api.GET("location", controllers.GetLocations)
	api.GET("location/:id", controllers.GetLocationByID)
	api.GET("location/:id/students", controllers.GetStudentsAtLocation)
	api.DELETE("location/:id", controllers.DeleteLocation)
	api.PATCH("location/:id", controllers.UpdateLocation)
	api.GET("location/:id/visits", controllers.VisitedLocationToday)
	api.POST("location/:id/logoutAll", controllers.LogoutAllStudentsAtLocation)

	api.POST("student", controllers.CreateStudent)
	api.POST("students", controllers.CreateStudents)
	api.GET("student", controllers.GetStudents)
	api.GET("student/:id", controllers.GetStudentByID)
	api.GET("student/:id/location", controllers.GetStudentLocation)
	api.DELETE("student/:id", controllers.DeleteStudent)
	api.PATCH("student/:id", controllers.UpdateStudent)
	api.POST("student/:id/logout", controllers.LogoutStudent)

	// Serve React frontend
	r.Use(static.Serve("/", static.LocalFile("frontend/build", false)))
	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/build/index.html")
	})

	return r.Run(addr)
}
