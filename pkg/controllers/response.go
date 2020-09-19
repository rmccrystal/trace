package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
		formattedErr := []rune(error.Error())
		formattedErr[0] = unicode.ToUpper(formattedErr[0])
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