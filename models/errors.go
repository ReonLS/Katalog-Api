package models

//implementing Error interface
type ErrorResponse struct{
	Message string `json:"message"`
	StatusCode int `json:"status_code"`
}

func (e *ErrorResponse) Error() string{
	return e.Message
}