package api

// Success should be called whenever a successful response needs to be sent
func Success(data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Success: true,
		Data:    data,
	}
}

type ResponseSuccess struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func Error(error error) ResponseError {
	return ResponseError{
		Success: false,
		Error:   error.Error(),
	}
}

type ResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
