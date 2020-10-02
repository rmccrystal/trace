package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"unicode"
)

// This file contains standardized responses depending on the success of a request
// If a request is successful

// Success should be called whenever a successful response needs to be sent.
// It will abort and reply with the following JSON:
// {
//   "success": true,
//   "data" <data>
// }
// If the data is nil, it will be omitted from the JSON
func Success(c *gin.Context, code int, data interface{}) {
	json := struct{
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Success: true,
		Data:    data,
	}
	c.AbortWithStatusJSON(code, json)
}

// Error should be called when a request fails and an unsuccessful status must be
// sent back to the client. It will abort the request and send a json response:
// {
//   "success": false,
//   "error" <error>
// }
// If the error is nil, the error element will be omitted from the JSON.
// The first letter of error will automatically be capitalized
func Error(c *gin.Context, code int, error error) {
	var formattedErr string
	if error.Error() != "" {
		runeErr := []rune(error.Error())
		runeErr[0] = unicode.ToUpper(runeErr[0])
		formattedErr = string(runeErr)
	} else {
		formattedErr = ""
	}

	json := struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}{
		Success: false,
		Error:   string(formattedErr),
	}
	c.AbortWithStatusJSON(code, json)
}

// Errorf return the same thing as Error except it formats the arguments
func Errorf(c *gin.Context, code int, format string, args ...interface{}) {
	Error(c, code, fmt.Errorf(format, args...))
}

// GetIDParam gets the :id parameter from a request and converts it into an object ID.
// If there is an error parsing the object ID, it will send an error response and the bool will be false.
// If the success bool is false, the caller should return from the request
func GetIDParam(c *gin.Context) (id primitive.ObjectID, success bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "could not parse object id: %s", err)
		return primitive.ObjectID{}, false
	}

	return id, true
}

// BindJSON calls gin.Context.BindJSON and responds with an error if it is unsuccessful.
// If the bool returned is false, the caller should return
func BindJSON(c *gin.Context, obj interface{}) bool {
	err := c.BindJSON(obj)
	if err != nil {
		Errorf(c, http.StatusUnprocessableEntity, "failed to parse request body: %s", err)
		return false
	}
	return true
}