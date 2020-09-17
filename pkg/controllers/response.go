package controllers

// This file contains standardized responses depending on the success of a request
// If a request is successful

// Success should be called whenever a successful response needs to be sent.
// It will return an interface{} which, when serialized, will turn into the following JSON:
// {
//   "success": true,
//   "data" <data>
// }
// If the data is nil, it will be omitted from the JSON
func Success(data interface{}) interface{} {
	return struct{
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Success: true,
		Data:    data,
	}
}

// Error should be called when a request fails and an unsuccessful status must be
// sent back to the client. It will return an interface which will serialize into:
// {
//   "success": false,
//   "error" <error>
// }
// If the error is nil, the error element will be omitted from the JSON
func Error(error error) interface{} {
	return struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}{
		Success: false,
		Error:   error.Error(),
	}
}
